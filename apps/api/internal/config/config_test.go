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

func establecerVariables(t *testing.T, entorno string, jwks string) {
	t.Helper()

	t.Setenv("API_PORT", "5171")
	t.Setenv("API_ENV", entorno)
	t.Setenv("DATABASE_URL", "postgres://pingeso@localhost:5433/pingeso?sslmode=disable")
	t.Setenv("AUTH_JWKS_URL", jwks)
	t.Setenv("FRONTEND_URL", "http://localhost:5170")
	t.Setenv("TZ", "America/Santiago")
}
