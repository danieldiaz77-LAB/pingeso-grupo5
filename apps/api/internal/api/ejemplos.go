package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
)

type Ejemplo struct {
	ID          string    `json:"id" format:"uuid" doc:"Identificador único del ejemplo"`
	Nombre      string    `json:"nombre" minLength:"1" maxLength:"120" doc:"Nombre del ejemplo"`
	Descripcion string    `json:"descripcion" maxLength:"500" doc:"Descripción del ejemplo"`
	CreadoEn    time.Time `json:"creadoEn" doc:"Fecha de creación"`
}

type ListarEjemplosEntrada struct{}

type ListarEjemplosSalida struct {
	Body struct {
		Ejemplos []Ejemplo `json:"ejemplos" doc:"Listado de ejemplos"`
	}
}

type CrearEjemploEntrada struct {
	Body struct {
		Nombre      string `json:"nombre" minLength:"1" maxLength:"120" doc:"Nombre del ejemplo"`
		Descripcion string `json:"descripcion" maxLength:"500" required:"false" doc:"Descripción del ejemplo"`
	}
}

type CrearEjemploSalida struct {
	Body Ejemplo
}

func registrarEjemplos(api huma.API, svc *service.EjemploServicio) {
	huma.Register(api, huma.Operation{
		OperationID: "listarEjemplos",
		Method:      http.MethodGet,
		Path:        "/api/v1/ejemplos",
		Summary:     "Lista los ejemplos disponibles",
		Tags:        []string{"Ejemplos"},
	}, func(ctx context.Context, _ *ListarEjemplosEntrada) (*ListarEjemplosSalida, error) {
		ejemplos, err := svc.Listar(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("no se pudo listar los ejemplos")
		}

		salida := &ListarEjemplosSalida{}
		salida.Body.Ejemplos = make([]Ejemplo, 0, len(ejemplos))
		for _, ejemplo := range ejemplos {
			salida.Body.Ejemplos = append(salida.Body.Ejemplos, aEjemplo(ejemplo))
		}
		return salida, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "crearEjemplo",
		Method:        http.MethodPost,
		Path:          "/api/v1/ejemplos",
		Summary:       "Crea un ejemplo",
		Tags:          []string{"Ejemplos"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, entrada *CrearEjemploEntrada) (*CrearEjemploSalida, error) {
		ejemplo, err := svc.Crear(ctx, entrada.Body.Nombre, entrada.Body.Descripcion)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrNombreRequerido),
				errors.Is(err, service.ErrNombreMuyLargo),
				errors.Is(err, service.ErrDescripcionMuyLarga):
				return nil, huma.Error400BadRequest(err.Error())
			default:
				return nil, huma.Error500InternalServerError("no se pudo crear el ejemplo")
			}
		}

		return &CrearEjemploSalida{Body: aEjemplo(ejemplo)}, nil
	})
}

func aEjemplo(ejemplo service.Ejemplo) Ejemplo {
	return Ejemplo{
		ID:          ejemplo.ID,
		Nombre:      ejemplo.Nombre,
		Descripcion: ejemplo.Descripcion,
		CreadoEn:    ejemplo.CreadoEn,
	}
}
