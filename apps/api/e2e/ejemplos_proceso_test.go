//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestEjemplosBinarioReal(t *testing.T) {
	base := baseDePrueba(t)
	cerrar := iniciar(t, "test", base)
	defer cerrar()
	url := "http://127.0.0.1:8080/api/v1/ejemplos"
	estado, cuerpo := solicitar(t, url, http.MethodPost, `{"nombre":"  Ejemplo integrado  ","descripcion":"  prueba  "}`)
	if estado != http.StatusCreated {
		t.Fatalf("crear: %d %s", estado, cuerpo)
	}
	var creado struct {
		ID          string `json:"id"`
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
	}
	if err := json.Unmarshal(cuerpo, &creado); err != nil {
		t.Fatal(err)
	}
	if creado.ID == "" || creado.Nombre != "Ejemplo integrado" || creado.Descripcion != "prueba" {
		t.Fatalf("normalización incorrecta: %+v", creado)
	}
	estado, cuerpo = solicitar(t, url, http.MethodGet, "")
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
	if !encontrado {
		t.Fatal("el ejemplo no persistió en PostgreSQL")
	}
	for _, caso := range []struct {
		entrada string
		estado  int
	}{
		{`{"nombre":"   "}`, 400},
		{`{}`, 422},
		{fmt.Sprintf(`{"nombre":"%s"}`, strings.Repeat("a", 121)), 422},
		{fmt.Sprintf(`{"nombre":"válido","descripcion":"%s"}`, strings.Repeat("a", 501)), 422},
	} {
		estado, cuerpo = solicitar(t, url, http.MethodPost, caso.entrada)
		if estado != caso.estado {
			t.Fatalf("rechazo esperado %d, recibido %d: %s", caso.estado, estado, cuerpo)
		}
	}
	consultar(t, "/readyz", 200)
	informe := "# Ejemplos de plantilla\n\n- Crear: 201; espacios normalizados\n- Listar: 200; registro persistido\n- Nombre en blanco: 400\n- Nombre ausente o demasiado largo: 422\n- Descripción demasiado larga: 422\n- PostgreSQL disponible: readyz 200\n"
	golden, err := os.ReadFile("testdata/ejemplos-proceso.golden.md")
	if err != nil {
		t.Fatal(err)
	}
	if informe != string(golden) {
		t.Fatalf("golden distinto:\n%s", informe)
	}
}
