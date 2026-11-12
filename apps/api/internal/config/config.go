package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Puerto      int
	Entorno     string
	DatabaseURL string
	AuthJWKSURL string
	FrontendURL string
	ZonaHoraria string
}

func Cargar() (Config, error) {
	puertoTexto := os.Getenv("API_PORT")
	if puertoTexto == "" {
		return Config{}, fmt.Errorf("falta la variable obligatoria API_PORT")
	}

	puerto, err := strconv.Atoi(puertoTexto)
	if err != nil || puerto < 1 || puerto > 65535 {
		return Config{}, fmt.Errorf("API_PORT inválido: %q", puertoTexto)
	}

	entorno := os.Getenv("API_ENV")
	if entorno == "" {
		return Config{}, fmt.Errorf("falta la variable obligatoria API_ENV")
	}
	if entorno != "dev" && entorno != "test" && entorno != "prod" {
		return Config{}, fmt.Errorf("API_ENV inválido: %q", entorno)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("falta la variable obligatoria DATABASE_URL")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		return Config{}, fmt.Errorf("falta la variable obligatoria FRONTEND_URL")
	}

	zonaHoraria := os.Getenv("TZ")
	if zonaHoraria == "" {
		return Config{}, fmt.Errorf("falta la variable obligatoria TZ")
	}

	authJWKSURL := os.Getenv("AUTH_JWKS_URL")
	if entorno == "prod" && authJWKSURL == "" {
		return Config{}, fmt.Errorf("falta la variable obligatoria AUTH_JWKS_URL en producción")
	}

	return Config{
		Puerto:      puerto,
		Entorno:     entorno,
		DatabaseURL: databaseURL,
		AuthJWKSURL: authJWKSURL,
		FrontendURL: frontendURL,
		ZonaHoraria: zonaHoraria,
	}, nil
}
