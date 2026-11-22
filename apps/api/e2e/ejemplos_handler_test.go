//go:build e2e

package e2e

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/api"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/repository"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/store"
)

func TestListarEjemplosContraLaBaseDeDatos(t *testing.T) {
	servidor := levantarEntorno(t)

	estado, cuerpo := solicitar(t, servidor.URL+"/api/v1/ejemplos", http.MethodGet, "")
	if estado != http.StatusOK {
		t.Fatalf("estado inesperado: %d, cuerpo: %s", estado, cuerpo)
	}

	golden, err := os.ReadFile("testdata/listar_ejemplos.golden.json")
	if err != nil {
		t.Fatalf("no se pudo leer el golden: %v", err)
	}

	var esperado any
	if err := json.Unmarshal(golden, &esperado); err != nil {
		t.Fatalf("golden inválido: %v", err)
	}
	var obtenido any
	if err := json.Unmarshal(cuerpo, &obtenido); err != nil {
		t.Fatalf("respuesta inválida: %v", err)
	}
	if objeto, ok := obtenido.(map[string]any); ok {
		delete(objeto, "$schema") // Metadato Huma dependiente de la URL de httptest.
	}
	if !reflect.DeepEqual(esperado, obtenido) {
		t.Fatalf("la respuesta no coincide con el golden\nobtenido: %s\nesperado: %s", cuerpo, golden)
	}
}

func TestCrearEjemploYVerificarloEnElListado(t *testing.T) {
	servidor := levantarEntorno(t)

	estado, cuerpo := solicitar(t, servidor.URL+"/api/v1/ejemplos", http.MethodPost, `{"nombre":"Cuarto ejemplo","descripcion":"Creado por la prueba e2e."}`)
	if estado != http.StatusCreated {
		t.Fatalf("estado inesperado al crear: %d, cuerpo: %s", estado, cuerpo)
	}

	var creado struct {
		ID          string `json:"id"`
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
	}
	if err := json.Unmarshal(cuerpo, &creado); err != nil {
		t.Fatalf("respuesta inválida al crear: %v", err)
	}
	if creado.ID == "" {
		t.Fatal("el ejemplo creado no tiene id")
	}
	if creado.Nombre != "Cuarto ejemplo" {
		t.Fatalf("nombre inesperado: %q", creado.Nombre)
	}

	estado, cuerpo = solicitar(t, servidor.URL+"/api/v1/ejemplos", http.MethodGet, "")
	if estado != http.StatusOK {
		t.Fatalf("estado inesperado al listar: %d, cuerpo: %s", estado, cuerpo)
	}

	var listado struct {
		Ejemplos []struct {
			ID     string `json:"id"`
			Nombre string `json:"nombre"`
		} `json:"ejemplos"`
	}
	if err := json.Unmarshal(cuerpo, &listado); err != nil {
		t.Fatalf("respuesta inválida al listar: %v", err)
	}
	if len(listado.Ejemplos) != 4 {
		t.Fatalf("se esperaban 4 ejemplos, se obtuvieron %d", len(listado.Ejemplos))
	}
	if listado.Ejemplos[3].ID != creado.ID {
		t.Fatalf("el último ejemplo no es el creado: %s", listado.Ejemplos[3].ID)
	}
}

func TestCrearEjemploRechazaNombreEnBlanco(t *testing.T) {
	servidor := levantarEntorno(t)

	estado, cuerpo := solicitar(t, servidor.URL+"/api/v1/ejemplos", http.MethodPost, `{"nombre":"   "}`)
	if estado != http.StatusBadRequest {
		t.Fatalf("estado inesperado: %d, cuerpo: %s", estado, cuerpo)
	}
}

func TestSaludYDisponibilidad(t *testing.T) {
	servidor := levantarEntorno(t)

	for _, ruta := range []string{"/healthz", "/readyz"} {
		estado, cuerpo := solicitar(t, servidor.URL+ruta, http.MethodGet, "")
		if estado != http.StatusOK {
			t.Fatalf("%s devolvió %d: %s", ruta, estado, cuerpo)
		}
	}
}

func baseDePrueba(t *testing.T) string {
	t.Helper()
	valor := os.Getenv("TEST_DATABASE_URL")
	nombre := os.Getenv("TEST_DATABASE_NAME")
	u, err := url.Parse(valor)
	if err != nil || valor == "" || !strings.HasPrefix(nombre, "pingeso_integracion_") || u.Path != "/"+nombre {
		t.Fatal("requiere TEST_DATABASE_URL y TEST_DATABASE_NAME de una base efímera pingeso_integracion_*")
	}
	return valor
}

func levantarEntorno(t *testing.T) *httptest.Server {
	t.Helper()
	anterior := time.Local
	time.Local = time.UTC
	t.Cleanup(func() { time.Local = anterior })
	ctx, cancelar := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelar()
	admin, err := pgxpool.New(ctx, baseDePrueba(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(admin.Close)
	if err := admin.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	var version string
	if err := admin.QueryRow(ctx, "SHOW server_version").Scan(&version); err != nil || !strings.HasPrefix(version, "17.11") {
		t.Fatalf("PostgreSQL debe ser 17.11: %q, %v", version, err)
	}
	var aleatorio [8]byte
	if _, err := rand.Read(aleatorio[:]); err != nil {
		t.Fatal(err)
	}
	esquema := "ejemplos_" + hex.EncodeToString(aleatorio[:])
	identificador := pgx.Identifier{esquema}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+identificador); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cierre, fin := context.WithTimeout(context.Background(), 10*time.Second)
		defer fin()
		if _, err := admin.Exec(cierre, "DROP SCHEMA "+identificador+" CASCADE"); err != nil {
			t.Error(err)
		}
	})
	urlEsquema, err := url.Parse(baseDePrueba(t))
	if err != nil {
		t.Fatal(err)
	}
	parametros := urlEsquema.Query()
	parametros.Set("search_path", esquema)
	parametros.Set("timezone", "UTC")
	urlEsquema.RawQuery = parametros.Encode()
	cfg, err := pgxpool.ParseConfig(urlEsquema.String())
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = esquema
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	// Goose y todas las conexiones de la API usan la misma URL con search_path.
	migrador := os.Getenv("GOOSE_E2E_BINARY")
	if migrador == "" {
		t.Fatal("falta GOOSE_E2E_BINARY; ejecutar pnpm test:e2e:api")
	}
	migracion := exec.CommandContext(ctx, migrador, "-dir", "../db/migrations", "up")
	for _, variable := range os.Environ() {
		clave, _, _ := strings.Cut(variable, "=")
		if clave != "GOOSE_DBSTRING" && clave != "GOOSE_DRIVER" && clave != "GOOSE_TABLE" {
			migracion.Env = append(migracion.Env, variable)
		}
	}
	migracion.Env = append(migracion.Env, "GOOSE_DRIVER=postgres", "GOOSE_DBSTRING="+urlEsquema.String(), "GOOSE_TABLE=goose_db_version")
	if err := migracion.Run(); err != nil {
		t.Fatalf("migración de esquema falló: %v", err)
	}
	var actual string
	if err := pool.QueryRow(ctx, "SELECT current_schema()").Scan(&actual); err != nil || actual != esquema {
		t.Fatalf("conexión fuera del esquema de prueba: %q, %v", actual, err)
	}
	var tablas int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name IN ('ejemplos', 'goose_db_version')", esquema).Scan(&tablas); err != nil || tablas != 2 {
		t.Fatalf("migración fuera del esquema: %d tablas, %v", tablas, err)
	}
	seed, err := os.ReadFile("../db/seed/seed.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, string(seed)); err != nil {
		t.Fatal(err)
	}
	svc := service.Nuevo(store.Nuevo(repository.New(pool)))
	handler, _ := api.Nueva(pool, "test", svc)
	servidor := httptest.NewServer(handler)
	t.Cleanup(servidor.Close)
	return servidor
}

func solicitar(t *testing.T, url string, metodo string, cuerpo string) (int, []byte) {
	t.Helper()

	var lector io.Reader
	if cuerpo != "" {
		lector = strings.NewReader(cuerpo)
	}

	solicitud, err := http.NewRequest(metodo, url, lector)
	if err != nil {
		t.Fatalf("no se pudo crear la solicitud: %v", err)
	}
	if cuerpo != "" {
		solicitud.Header.Set("Content-Type", "application/json")
	}

	cliente := &http.Client{Timeout: 10 * time.Second}
	respuesta, err := cliente.Do(solicitud)
	if err != nil {
		t.Fatalf("la solicitud falló: %v", err)
	}
	defer respuesta.Body.Close()

	contenido, err := io.ReadAll(respuesta.Body)
	if err != nil {
		t.Fatalf("no se pudo leer la respuesta: %v", err)
	}
	return respuesta.StatusCode, contenido
}

// Dos fixtures simultáneos deben compartir solo la instancia de PostgreSQL.
func TestAislamientoDeEsquemas(t *testing.T) {
	primero := levantarEntorno(t)
	segundo := levantarEntorno(t)
	estado, cuerpo := solicitar(t, primero.URL+"/api/v1/ejemplos", http.MethodPost, `{"nombre":"Solo primer esquema"}`)
	if estado != http.StatusCreated {
		t.Fatalf("crear: %d %s", estado, cuerpo)
	}
	var creado struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(cuerpo, &creado); err != nil || creado.ID == "" {
		t.Fatal("respuesta de creación inválida", err)
	}
	for _, caso := range []struct {
		servidor *httptest.Server
		esperado int
		contiene bool
	}{
		{primero, 4, true}, {segundo, 3, false},
	} {
		estado, cuerpo = solicitar(t, caso.servidor.URL+"/api/v1/ejemplos", http.MethodGet, "")
		if estado != http.StatusOK {
			t.Fatalf("listar: %d %s", estado, cuerpo)
		}
		var listado struct {
			Ejemplos []struct {
				ID string `json:"id"`
			} `json:"ejemplos"`
		}
		if err := json.Unmarshal(cuerpo, &listado); err != nil {
			t.Fatal(err)
		}
		encontrado := false
		for _, fila := range listado.Ejemplos {
			if fila.ID == creado.ID {
				encontrado = true
			}
		}
		if len(listado.Ejemplos) != caso.esperado || encontrado != caso.contiene {
			t.Fatal("datos compartidos entre esquemas")
		}
	}
}
