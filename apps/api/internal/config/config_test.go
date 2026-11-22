package config

import (
	"strings"
	"testing"
)

func TestCargarRechazaProduccionSinJWKS(t *testing.T) {
	establecerVariables(t, "prod", "")

	_, err := Cargar()
	if err == nil || !strings.Contains(err.Error(), "AUTH_JWKS_URL") {
		t.Fatalf("se esperaba un error por AUTH_JWKS_URL, se obtuvo: %v", err)
	}
}

func TestCargarPermiteDesarrolloYSinJWKS(t *testing.T) {
	establecerVariables(t, "dev", "")

	cfg, err := Cargar()
	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}
	if cfg.AuthJWKSURL != "" {
		t.Fatalf("AUTH_JWKS_URL inesperado: %q", cfg.AuthJWKSURL)
	}
}

func TestCargarProduccionConJWKS(t *testing.T) {
	establecerVariables(t, "prod", "https://identidad.interna/.well-known/jwks.json")

	cfg, err := Cargar()
	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}
	if cfg.Entorno != "prod" || cfg.AuthJWKSURL == "" || cfg.ZonaHoraria != "America/Santiago" {
		t.Fatalf("configuración incompleta: %+v", cfg)
	}
}

func TestCargarRechazaValoresInvalidos(t *testing.T) {
	for _, caso := range []struct{ clave, valor, mensaje string }{
		{"API_PORT", "", "falta API_PORT"},
		{"API_PORT", "5171", "API_PORT debe ser 8080"},
		{"API_ENV", "otro", "API_ENV inválido"},
		{"DATABASE_URL", "", "falta DATABASE_URL"},
		{"FRONTEND_URL", "file:///tmp", "FRONTEND_URL inválido"},
		{"FRONTEND_URL", "://", "FRONTEND_URL inválido"},
		{"TZ", "Zona/Inexistente", "TZ inválido"},
		{"AUTH_JWKS_URL", "replace-with-identity-jwks-url", "AUTH_JWKS_URL inválido"},
		{"AUTH_JWKS_URL", "ftp://identidad/jwks.json", "AUTH_JWKS_URL inválido"},
	} {
		t.Run(caso.clave+"_"+caso.valor, func(t *testing.T) {
			establecerVariables(t, "test", "")
			t.Setenv(caso.clave, caso.valor)
			_, err := Cargar()
			if err == nil || err.Error() != caso.mensaje {
				t.Fatalf("esperado %q, obtenido %v", caso.mensaje, err)
			}
		})
	}
}

func establecerVariables(t *testing.T, entorno string, jwks string) {
	t.Helper()

	t.Setenv("API_PORT", "8080")
	t.Setenv("API_ENV", entorno)
	t.Setenv("DATABASE_URL", "postgres://pingeso@localhost:5433/pingeso?sslmode=disable")
	t.Setenv("AUTH_JWKS_URL", jwks)
	t.Setenv("FRONTEND_URL", "http://localhost:5170")
	t.Setenv("TZ", "America/Santiago")
}
