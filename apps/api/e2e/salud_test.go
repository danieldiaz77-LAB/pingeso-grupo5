//go:build e2e

package e2e

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func iniciar(t *testing.T, ambiente, base string) func() {
	t.Helper()
	cmd := exec.Command(os.Getenv("API_E2E_BINARY"))
	cmd.Env = append(os.Environ(), "API_ENV="+ambiente, "API_PORT=8080", "DATABASE_URL="+base, "FRONTEND_URL=http://127.0.0.1:8080", "TZ=America/Santiago")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	cerrar := func() {
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			t.Error(err)
			_ = cmd.Process.Kill()
		}
		terminado := make(chan error, 1)
		go func() { terminado <- cmd.Wait() }()
		select {
		case err := <-terminado:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(15 * time.Second):
			_ = cmd.Process.Kill()
			t.Error("timeout de cierre")
		}
	}
	for i := 0; i < 100; i++ {
		cliente := http.Client{Timeout: time.Second}
		r, err := cliente.Get("http://127.0.0.1:8080/healthz")
		if err == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				return cerrar
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	cerrar()
	t.Fatal("API no inició")
	return nil
}
func consultar(t *testing.T, ruta string, codigo int) string {
	t.Helper()
	cliente := http.Client{Timeout: 7 * time.Second}
	r, err := cliente.Get("http://127.0.0.1:8080" + ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	datos, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != codigo {
		t.Fatalf("%s: esperado %d, recibido %d: %s", ruta, codigo, r.StatusCode, datos)
	}
	return string(datos)
}
func TestSaludReal(t *testing.T) {
	base := os.Getenv("TEST_DATABASE_URL")
	if base == "" {
		t.Fatal("TEST_DATABASE_URL obligatorio: requiere PostgreSQL real")
	}
	cerrar := iniciar(t, "test", base)
	consultar(t, "/readyz", 200)
	consultar(t, "/api/openapi.yaml", 200)
	var grupo sync.WaitGroup
	for i := 0; i < 8; i++ {
		grupo.Add(1)
		go func() { defer grupo.Done(); consultar(t, "/readyz", 200) }()
	}
	grupo.Wait()
	cerrar()
	cerrar = iniciar(t, "prod", "postgres://prueba:prueba@127.0.0.1:1/prueba?connect_timeout=1")
	consultar(t, "/healthz", 200)
	cuerpo := consultar(t, "/readyz", 503)
	consultar(t, "/api/openapi.yaml", 404)
	consultar(t, "/api/docs", 404)
	if strings.Contains(cuerpo, "$schema") || strings.Contains(cuerpo, "postgres://") || strings.Contains(cuerpo, "127.0.0.1") || !strings.Contains(cuerpo, "Base de datos no disponible") {
		t.Error("error público no sanitizado")
	}
	cerrar()
	informe := fmt.Sprintln("# Salud técnica") + "\n- PostgreSQL disponible: readyz 200\n- Ocho consultas concurrentes: 200\n- PostgreSQL inaccesible: readyz 503\n- Proceso vivo sin PostgreSQL: healthz 200\n- Contrato y docs en prod: 404\n"
	golden, err := os.ReadFile(filepath.Join("testdata", "salud.golden.md"))
	if err != nil {
		t.Fatal(err)
	}
	if informe != string(golden) {
		t.Fatalf("golden distinto:\n%s", informe)
	}
}

func TestServidorWebReal(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("inicio"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "assets"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("codigo"), 0600); err != nil {
		t.Fatal(err)
	}
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("backend")) }))
	defer backend.Close()
	cmd := exec.Command(os.Getenv("WEB_E2E_BINARY"))
	cmd.Env = append(os.Environ(), "BACKEND_URL="+backend.URL, "STATIC_DIR="+dir)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		terminado := make(chan error, 1)
		go func() { terminado <- cmd.Wait() }()
		select {
		case err := <-terminado:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(15 * time.Second):
			_ = cmd.Process.Kill()
			t.Error("web no terminó")
		}
	}()
	listo := false
	for i := 0; i < 100; i++ {
		cliente := http.Client{Timeout: time.Second}
		r, err := cliente.Get("http://127.0.0.1:8080/healthz")
		if err == nil {
			r.Body.Close()
			listo = r.StatusCode == 200
			if listo {
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !listo {
		t.Fatal("web no inició")
	}
	for ruta, esperado := range map[string]string{"/ruta-spa": "inicio", "/api/prueba": "backend", "/assets/app.js": "codigo"} {
		cuerpo := consultar(t, ruta, 200)
		if strings.TrimSpace(cuerpo) != esperado {
			t.Fatalf("%s: %s", ruta, cuerpo)
		}
	}
	for _, ruta := range []string{"/../privado", "/%2e%2e/privado", "/assets/%2e%2e/%2e%2e/privado", "/ruta%5cprivada"} {
		consultar(t, ruta, 400)
	}
	consultar(t, "/assets/ausente.js", 404)
	consultar(t, "/assets/%252e%252e/privado", 404)
	for metodo, codigo := range map[string]int{"HEAD": 200, "POST": 405} {
		peticion, _ := http.NewRequest(metodo, "http://127.0.0.1:8080/assets/app.js", nil)
		respuesta, err := http.DefaultClient.Do(peticion)
		if err != nil {
			t.Fatal(err)
		}
		cuerpo, err := io.ReadAll(respuesta.Body)
		respuesta.Body.Close()
		if err != nil || respuesta.StatusCode != codigo {
			t.Fatalf("%s: %v %d", metodo, err, respuesta.StatusCode)
		}
		if metodo == "HEAD" && len(cuerpo) != 0 {
			t.Fatal("HEAD incluyó cuerpo")
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("reemplazado"), 0600); err != nil {
		t.Fatal(err)
	}
	if consultar(t, "/assets/app.js", 200) != "codigo" {
		t.Fatal("el catálogo cambió después del arranque")
	}
	backend.Close()
	consultar(t, "/api/falla", 502)
}
