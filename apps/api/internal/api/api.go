package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/config"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
)

type VerificadorBD interface {
	Ping(ctx context.Context) error
}

func NuevaHuma(cfg config.Config, svc *service.EjemploServicio) (huma.API, *http.ServeMux) {
	mux := http.NewServeMux()

	configHuma := huma.DefaultConfig("Pingeso Grupo 5 API", "0.1.0")
	configHuma.Info.Description = "API de demostración de la plantilla técnica del Grupo 5."
	// El transformador por defecto agrega "$schema" y Link con la URL del
	// request. El contrato del servicio no expone esos campos.
	configHuma.CreateHooks = nil

	if cfg.Entorno == "prod" {
		configHuma.OpenAPIPath = ""
		configHuma.DocsPath = ""
		configHuma.SchemasPath = ""
	} else {
		configHuma.OpenAPIPath = "/api/openapi"
		configHuma.DocsPath = "/api/docs"
		configHuma.SchemasPath = "/api/schemas"
	}

	api := humago.New(mux, configHuma)
	registrarEjemplos(api, svc)

	return api, mux
}

func Construir(
	cfg config.Config,
	svc *service.EjemploServicio,
	salud VerificadorBD,
	logger *slog.Logger,
	version string,
) (http.Handler, huma.API) {
	api, mux := NuevaHuma(cfg, svc)
	registrarSalud(mux, salud)
	return aplicarMiddleware(mux, cfg, logger, version), api
}
