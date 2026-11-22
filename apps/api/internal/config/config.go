package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"
)

// Config reúne las variables de entorno del DAS §5. El puerto interno es
// siempre 8080, por eso no se guarda.
type Config struct {
	Entorno     string
	DatabaseURL string
	AuthJWKSURL string
	FrontendURL string
	ZonaHoraria string
}

// Cargar valida la configuración al arrancar; cualquier error debe terminar el
// proceso. Los mensajes no incluyen valores para no filtrar credenciales.
func Cargar() (Config, error) {
	for _, clave := range []string{"API_PORT", "API_ENV", "DATABASE_URL", "FRONTEND_URL", "TZ"} {
		if os.Getenv(clave) == "" {
			return Config{}, fmt.Errorf("falta %s", clave)
		}
	}
	if os.Getenv("API_PORT") != "8080" {
		return Config{}, errors.New("API_PORT debe ser 8080")
	}
	cfg := Config{
		Entorno:     os.Getenv("API_ENV"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AuthJWKSURL: os.Getenv("AUTH_JWKS_URL"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
		ZonaHoraria: os.Getenv("TZ"),
	}
	if cfg.Entorno != "dev" && cfg.Entorno != "test" && cfg.Entorno != "prod" {
		return Config{}, errors.New("API_ENV inválido")
	}
	if !urlHTTP(cfg.FrontendURL) {
		return Config{}, errors.New("FRONTEND_URL inválido")
	}
	if _, err := time.LoadLocation(cfg.ZonaHoraria); err != nil {
		return Config{}, errors.New("TZ inválido")
	}
	if cfg.AuthJWKSURL == "" {
		if cfg.Entorno == "prod" {
			return Config{}, errors.New("falta AUTH_JWKS_URL en prod")
		}
	} else if !urlHTTP(cfg.AuthJWKSURL) {
		return Config{}, errors.New("AUTH_JWKS_URL inválido")
	}
	return cfg, nil
}

func urlHTTP(valor string) bool {
	destino, err := url.Parse(valor)
	return err == nil && destino.Host != "" && (destino.Scheme == "http" || destino.Scheme == "https")
}
