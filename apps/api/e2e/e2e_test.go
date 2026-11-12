package e2e

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/api"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/config"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/repository"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/store"
)

const (
	imagenPostgres = "postgres:17.11@sha256:ff87138949ad9a46f32ff23bd6852b52ce024a13dffe9b4087be98292cdd3b4f"
	directorioSQL  = "../db"
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

func levantarEntorno(t *testing.T) *httptest.Server {
	t.Helper()

	// La serialización de timestamptz depende de time.Local. La prueba fija UTC
	// para que el golden sea determinista en cualquier equipo.
	time.Local = time.UTC

	urlBaseDatos := urlDeBaseDeDatos(t)

	dbSQL, err := sql.Open("pgx", urlBaseDatos)
	if err != nil {
		t.Fatalf("no se pudo abrir la conexión SQL: %v", err)
	}
	t.Cleanup(func() { _ = dbSQL.Close() })

	var versionPostgres string
	if err := dbSQL.QueryRow("SHOW server_version").Scan(&versionPostgres); err != nil {
		t.Fatalf("no se pudo leer la versión de PostgreSQL: %v", err)
	}
	if !strings.HasPrefix(versionPostgres, "17.11") {
		t.Fatalf("PostgreSQL %s no coincide con la versión exacta 17.11 del DAS", versionPostgres)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("no se pudo configurar goose: %v", err)
	}
	if err := goose.Up(dbSQL, directorioSQL+"/migrations"); err != nil {
		t.Fatalf("no se pudieron aplicar las migraciones: %v", err)
	}

	if _, err := dbSQL.Exec("TRUNCATE TABLE ejemplos"); err != nil {
		t.Fatalf("no se pudo limpiar la tabla: %v", err)
	}

	seed, err := os.ReadFile(directorioSQL + "/seed/seed.sql")
	if err != nil {
		t.Fatalf("no se pudo leer el seed: %v", err)
	}
	if _, err := dbSQL.Exec(string(seed)); err != nil {
		t.Fatalf("no se pudo aplicar el seed: %v", err)
	}

	ctx, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()

	pool, err := pgxpool.New(ctx, urlBaseDatos)
	if err != nil {
		t.Fatalf("no se pudo crear el pool: %v", err)
	}
	t.Cleanup(pool.Close)

	cfg := config.Config{
		Entorno:     "test",
		Puerto:      0,
		DatabaseURL: urlBaseDatos,
		FrontendURL: "http://localhost:5170",
		ZonaHoraria: "UTC",
	}
	servicio := service.Nuevo(store.Nuevo(repository.New(pool)))
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler, _ := api.Construir(cfg, servicio, pool, logger, "test")

	servidor := httptest.NewServer(handler)
	t.Cleanup(servidor.Close)
	return servidor
}

func urlDeBaseDeDatos(t *testing.T) string {
	t.Helper()

	if url := strings.TrimSpace(os.Getenv("TEST_DATABASE_URL")); url != "" {
		return url
	}

	ctx, cancelar := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelar()

	contenedor, err := postgres.Run(ctx, imagenPostgres,
		postgres.WithDatabase("pingeso"),
		postgres.WithUsername("pingeso"),
		postgres.WithPassword("pingeso"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Skipf("Docker no está disponible para la prueba e2e: %v", err)
	}
	t.Cleanup(func() {
		_ = testcontainers.TerminateContainer(contenedor)
	})

	url, err := contenedor.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("no se pudo obtener la cadena de conexión: %v", err)
	}
	return url
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
