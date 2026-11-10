//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func ejecutarBinario(t *testing.T, binario string, variables map[string]string, correcto bool, args ...string) string {
	t.Helper()
	ctx, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	ruta := os.Getenv(binario)
	if ruta == "" {
		t.Fatalf("falta %s", binario)
	}
	if _, err := os.Stat(ruta); err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, ruta, args...)
	for _, variable := range os.Environ() {
		clave, _, _ := strings.Cut(variable, "=")
		if _, reemplazada := variables[clave]; !reemplazada {
			cmd.Env = append(cmd.Env, variable)
		}
	}
	for clave, valor := range variables {
		cmd.Env = append(cmd.Env, clave+"="+valor)
	}
	salida, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("proceso no terminó: %s", salida)
	}
	var salidaProceso *exec.ExitError
	if err != nil && !errors.As(err, &salidaProceso) {
		t.Fatalf("no se pudo ejecutar el binario: %v", err)
	}
	if (err == nil) != correcto {
		t.Fatalf("éxito esperado=%v, error=%v: %s", correcto, err, salida)
	}
	return string(salida)
}
func configuracionAPI() map[string]string {
	return map[string]string{"API_PORT": "8080", "API_ENV": "test", "DATABASE_URL": os.Getenv("TEST_DATABASE_URL"), "FRONTEND_URL": "http://localhost:8080", "TZ": "America/Santiago"}
}
func TestRechazosDeArranque(t *testing.T) {
	for _, caso := range []struct{ clave, valor, mensaje string }{
		{"API_PORT", "", "falta API_PORT"}, {"API_ENV", "", "falta API_ENV"},
		{"DATABASE_URL", "", "falta DATABASE_URL"}, {"FRONTEND_URL", "", "falta FRONTEND_URL"}, {"TZ", "", "falta TZ"},
		{"API_PORT", "8090", "API_PORT debe ser"}, {"API_ENV", "otro", "API_ENV inválido"},
		{"FRONTEND_URL", "://", "FRONTEND_URL inválido"}, {"FRONTEND_URL", "file:///tmp", "FRONTEND_URL inválido"},
		{"TZ", "Zona/Inexistente", "TZ inválido"}, {"DATABASE_URL", "postgres://%", "DATABASE_URL inválido"},
	} {
		t.Run(caso.clave+"_"+caso.valor, func(t *testing.T) {
			variables := configuracionAPI()
			variables[caso.clave] = caso.valor
			salida := ejecutarBinario(t, "API_E2E_BINARY", variables, false)
			if !strings.Contains(salida, caso.mensaje) {
				t.Fatalf("error incorrecto: %s", salida)
			}
		})
	}
	ejecutarBinario(t, "WEB_E2E_BINARY", map[string]string{"BACKEND_URL": "invalido", "STATIC_DIR": ""}, false)
	ejecutarBinario(t, "WEB_E2E_BINARY", map[string]string{"BACKEND_URL": "http://localhost:8081", "STATIC_DIR": filepath.Join(t.TempDir(), "ausente")}, false)
}
func TestHealthchecksYConflictoDePuerto(t *testing.T) {
	for _, estado := range []int{200, 503} {
		listener, err := net.Listen("tcp", "127.0.0.1:8080")
		if err != nil {
			t.Fatal(err)
		}
		servidor := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(estado) }), ReadHeaderTimeout: time.Second}
		go servidor.Serve(listener)
		for _, binario := range []string{"API_E2E_BINARY", "WEB_E2E_BINARY"} {
			ejecutarBinario(t, binario, nil, estado == 200, "healthcheck")
		}
		if estado == 200 {
			salida := ejecutarBinario(t, "API_E2E_BINARY", configuracionAPI(), false)
			if !strings.Contains(salida, "no se pudo servir HTTP") {
				t.Fatal(salida)
			}
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("inicio"), 0600); err != nil {
				t.Fatal(err)
			}
			ejecutarBinario(t, "WEB_E2E_BINARY", map[string]string{"BACKEND_URL": "http://localhost:8081", "STATIC_DIR": dir}, false)
		}
		servidor.Close()
	}
	for _, binario := range []string{"API_E2E_BINARY", "WEB_E2E_BINARY"} {
		ejecutarBinario(t, binario, nil, false, "healthcheck")
	}
}
func TestContratoDelBinario(t *testing.T) {
	salida := ejecutarBinario(t, "OPENAPI_E2E_BINARY", nil, true)
	contrato, err := os.ReadFile("../../../contracts/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(salida), contrato) {
		t.Fatal("contrato generado distinto al versionado")
	}
}
func TestBundleSinIndice(t *testing.T) {
	ejecutarBinario(t, "WEB_E2E_BINARY", map[string]string{"BACKEND_URL": "http://localhost:8081", "STATIC_DIR": t.TempDir()}, false)
}
func TestBundleConEnlaceExterno(t *testing.T) {
	dir := t.TempDir()
	externo := filepath.Join(t.TempDir(), "privado.txt")
	for nombre, contenido := range map[string]string{filepath.Join(dir, "index.html"): "inicio", externo: "dato privado"} {
		if err := os.WriteFile(nombre, []byte(contenido), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(externo, filepath.Join(dir, "enlace.txt")); err != nil {
		t.Fatal(err)
	}
	ejecutarBinario(t, "WEB_E2E_BINARY", map[string]string{"BACKEND_URL": "http://localhost:8081", "STATIC_DIR": dir}, false)
}
