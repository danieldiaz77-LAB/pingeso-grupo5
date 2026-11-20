package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrNombreRequerido     = errors.New("el nombre es obligatorio")
	ErrNombreMuyLargo      = errors.New("el nombre supera los 120 caracteres")
	ErrDescripcionMuyLarga = errors.New("la descripción supera los 500 caracteres")
)

type Ejemplo struct {
	ID          string
	Nombre      string
	Descripcion string
	CreadoEn    time.Time
}

type AlmacenEjemplos interface {
	Listar(ctx context.Context) ([]repository.Ejemplo, error)
	Crear(ctx context.Context, nombre string, descripcion string) (repository.Ejemplo, error)
}

type EjemploServicio struct {
	almacen AlmacenEjemplos
}

func Nuevo(almacen AlmacenEjemplos) *EjemploServicio {
	return &EjemploServicio{almacen: almacen}
}

func (s *EjemploServicio) Listar(ctx context.Context) ([]Ejemplo, error) {
	filas, err := s.almacen.Listar(ctx)
	if err != nil {
		return nil, fmt.Errorf("listar ejemplos: %w", err)
	}

	ejemplos := make([]Ejemplo, 0, len(filas))
	for _, fila := range filas {
		ejemplos = append(ejemplos, desdeRepositorio(fila))
	}
	return ejemplos, nil
}

func (s *EjemploServicio) Crear(ctx context.Context, nombre string, descripcion string) (Ejemplo, error) {
	nombre = strings.TrimSpace(nombre)
	descripcion = strings.TrimSpace(descripcion)

	if nombre == "" {
		return Ejemplo{}, ErrNombreRequerido
	}
	if utf8.RuneCountInString(nombre) > 120 {
		return Ejemplo{}, ErrNombreMuyLargo
	}
	if utf8.RuneCountInString(descripcion) > 500 {
		return Ejemplo{}, ErrDescripcionMuyLarga
	}

	fila, err := s.almacen.Crear(ctx, nombre, descripcion)
	if err != nil {
		return Ejemplo{}, fmt.Errorf("crear ejemplo: %w", err)
	}
	return desdeRepositorio(fila), nil
}

func desdeRepositorio(fila repository.Ejemplo) Ejemplo {
	return Ejemplo{
		ID:          formatearUUID(fila.ID),
		Nombre:      fila.Nombre,
		Descripcion: fila.Descripcion,
		CreadoEn:    fila.CreadoEn.Time,
	}
}

func formatearUUID(valor pgtype.UUID) string {
	if !valor.Valid {
		return ""
	}
	bytes := valor.Bytes
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	)
}
