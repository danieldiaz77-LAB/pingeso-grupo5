package e2e

import (
	"context"
	"errors"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/api"
	"net/http/httptest"
	"testing"
)

type baseCaida struct{}

func (baseCaida) Ping(context.Context) error { return errors.New("detalle privado") }

// Prueba aislada de la protección; no reemplaza el e2e PostgreSQL.
func TestReadiness(t *testing.T) {
	handler, _ := api.Nueva(baseCaida{}, "prod")
	for ruta, esperado := range map[string]int{"/readyz": 503, "/healthz": 200, "/api/docs": 404, "/api/openapi.yaml": 404} {
		salida := httptest.NewRecorder()
		handler.ServeHTTP(salida, httptest.NewRequest("GET", ruta, nil))
		if salida.Code != esperado {
			t.Fatalf("%s: esperado %d, recibido %d", ruta, esperado, salida.Code)
		}
	}
}
