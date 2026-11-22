package store

import (
	"context"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/repository"
)

type Ejemplos struct {
	consultas *repository.Queries
}

func Nuevo(consultas *repository.Queries) *Ejemplos {
	return &Ejemplos{consultas: consultas}
}

func (s *Ejemplos) Listar(ctx context.Context) ([]repository.Ejemplo, error) {
	return s.consultas.ListarEjemplos(ctx)
}

func (s *Ejemplos) Crear(ctx context.Context, nombre string, descripcion string) (repository.Ejemplo, error) {
	return s.consultas.CrearEjemplo(ctx, repository.CrearEjemploParams{
		Nombre:      nombre,
		Descripcion: descripcion,
	})
}
