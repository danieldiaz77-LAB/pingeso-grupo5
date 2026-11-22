package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecuperarResponde500SinTraza(t *testing.T) {
	manejador := Recuperar(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(errors.New("detalle interno /ruta/secreta"))
	}))
	respuesta := httptest.NewRecorder()
	respuesta.Header().Set("X-Request-Id", "id-prueba")
	manejador.ServeHTTP(respuesta, httptest.NewRequest(http.MethodGet, "/api/v1/ejemplos", nil))
	if respuesta.Code != http.StatusInternalServerError {
		t.Fatalf("estado %d", respuesta.Code)
	}
	if respuesta.Header().Get("Content-Type") != "application/problem+json" {
		t.Fatalf("content-type %q", respuesta.Header().Get("Content-Type"))
	}
	cuerpo := respuesta.Body.String()
	if strings.Contains(cuerpo, "secreta") || strings.Contains(cuerpo, "goroutine") || !strings.Contains(cuerpo, `"status":500`) {
		t.Fatalf("cuerpo inesperado: %s", cuerpo)
	}
}

func TestRecuperarRelanzaAbort(t *testing.T) {
	manejador := Recuperar(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic(http.ErrAbortHandler) }))
	defer func() {
		if recuperado := recover(); recuperado != http.ErrAbortHandler {
			t.Fatalf("se esperaba ErrAbortHandler, se obtuvo %v", recuperado)
		}
	}()
	manejador.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}

func TestRecuperarSinPanicoNoInterfiere(t *testing.T) {
	manejador := Recuperar(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusAccepted) }))
	respuesta := httptest.NewRecorder()
	manejador.ServeHTTP(respuesta, httptest.NewRequest(http.MethodGet, "/", nil))
	if respuesta.Code != http.StatusAccepted {
		t.Fatalf("estado %d", respuesta.Code)
	}
}
