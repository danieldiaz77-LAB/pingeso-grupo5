package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type Comprobador interface{ Ping(context.Context) error }
type Salud struct {
	Body struct {
		Estado string `json:"estado" enum:"ok" example:"ok"`
	}
}

func Nueva(base Comprobador, ambiente string) (*http.ServeMux, huma.API) {
	mux := http.NewServeMux()
	config := huma.DefaultConfig("Base técnica PINGESO", "0.0.0")
	config.DocsPath = "/api/docs"
	config.OpenAPIPath = "/api/openapi"
	config.SchemasPath = "/api/schemas"
	if ambiente == "prod" {
		config.DocsPath = ""
		config.OpenAPIPath = ""
		config.SchemasPath = ""
		// No publicar enlaces a esquemas cuando sus rutas estan deshabilitadas.
		config.CreateHooks = nil
	}
	contrato := humago.New(mux, config)
	huma.Get(contrato, "/healthz", func(ctx context.Context, entrada *struct{}) (*Salud, error) {
		salida := &Salud{}
		salida.Body.Estado = "ok"
		return salida, nil
	})
	huma.Get(contrato, "/readyz", func(ctx context.Context, entrada *struct{}) (*Salud, error) {
		ctx, cancelar := context.WithTimeout(ctx, 5*time.Second)
		defer cancelar()
		if err := base.Ping(ctx); err != nil {
			return nil, huma.Error503ServiceUnavailable("Base de datos no disponible")
		}
		salida := &Salud{}
		salida.Body.Estado = "ok"
		return salida, nil
	})
	return mux, contrato
}
