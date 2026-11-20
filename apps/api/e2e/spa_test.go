package e2e

import (
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/spa"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSPA(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("pagina tecnica"), 0600); err != nil {
		t.Fatal(err)
	}
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(r.URL.Path)) }))
	defer backend.Close()
	h, err := spa.Nuevo(dir, backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	for ruta, cuerpo := range map[string]string{"/healthz": `{"estado":"ok"}`, "/ruta-local": "pagina tecnica", "/api/ejemplo": "/api/ejemplo"} {
		r := httptest.NewRecorder()
		h.ServeHTTP(r, httptest.NewRequest("GET", ruta, nil))
		datos, _ := io.ReadAll(r.Result().Body)
		if r.Code != 200 || string(datos) != cuerpo {
			t.Fatalf("%s: %d %s", ruta, r.Code, datos)
		}
		if r.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatal("falta nosniff")
		}
	}
	if _, err := spa.Nuevo(dir, "invalido"); err == nil {
		t.Fatal("aceptó backend inválido")
	}
}
