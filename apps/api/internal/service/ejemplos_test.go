package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/repository"
)

var uuidDePrueba = pgtype.UUID{
	Bytes: [16]byte{
		0x11, 0x11, 0x11, 0x11,
		0x11, 0x11, 0x41, 0x11,
		0x81, 0x11, 0x11, 0x11,
		0x11, 0x11, 0x11, 0x11,
	},
	Valid: true,
}

var fechaDePrueba = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

type almacenFalso struct {
	filas       []repository.Ejemplo
	err         error
	nombre      string
	descripcion string
}

func (a *almacenFalso) Listar(_ context.Context) ([]repository.Ejemplo, error) {
	return a.filas, a.err
}

func (a *almacenFalso) Crear(_ context.Context, nombre string, descripcion string) (repository.Ejemplo, error) {
	a.nombre = nombre
	a.descripcion = descripcion
	if a.err != nil {
		return repository.Ejemplo{}, a.err
	}
	return repository.Ejemplo{
		ID:          uuidDePrueba,
		Nombre:      nombre,
		Descripcion: descripcion,
		CreadoEn:    pgtype.Timestamptz{Time: fechaDePrueba, Valid: true},
	}, nil
}

func TestListarConvierteLosDatosDelRepositorio(t *testing.T) {
	almacen := &almacenFalso{
		filas: []repository.Ejemplo{{
			ID:          uuidDePrueba,
			Nombre:      "Primer ejemplo",
			Descripcion: "Descripción",
			CreadoEn:    pgtype.Timestamptz{Time: fechaDePrueba, Valid: true},
		}},
	}
	servicio := Nuevo(almacen)

	ejemplos, err := servicio.Listar(context.Background())
	if err != nil {
		t.Fatalf("Listar devolvió error: %v", err)
	}
	if len(ejemplos) != 1 {
		t.Fatalf("se esperaba 1 ejemplo, se obtuvo %d", len(ejemplos))
	}
	if ejemplos[0].ID != "11111111-1111-4111-8111-111111111111" {
		t.Fatalf("UUID inesperado: %q", ejemplos[0].ID)
	}
	if !ejemplos[0].CreadoEn.Equal(fechaDePrueba) {
		t.Fatalf("fecha inesperada: %v", ejemplos[0].CreadoEn)
	}
}

func TestListarPropagaElErrorDelAlmacen(t *testing.T) {
	servicio := Nuevo(&almacenFalso{err: errors.New("sin conexión")})

	_, err := servicio.Listar(context.Background())
	if err == nil {
		t.Fatal("se esperaba un error")
	}
	if !strings.Contains(err.Error(), "listar ejemplos") {
		t.Fatalf("error sin contexto: %v", err)
	}
}

func TestCrearRecortaEspaciosYDelega(t *testing.T) {
	almacen := &almacenFalso{}
	servicio := Nuevo(almacen)

	ejemplo, err := servicio.Crear(context.Background(), "  Nuevo  ", "  Descripción  ")
	if err != nil {
		t.Fatalf("Crear devolvió error: %v", err)
	}
	if almacen.nombre != "Nuevo" {
		t.Fatalf("nombre enviado al almacén: %q", almacen.nombre)
	}
	if almacen.descripcion != "Descripción" {
		t.Fatalf("descripción enviada al almacén: %q", almacen.descripcion)
	}
	if ejemplo.Nombre != "Nuevo" {
		t.Fatalf("nombre devuelto: %q", ejemplo.Nombre)
	}
}

func TestCrearRechazaNombreVacio(t *testing.T) {
	servicio := Nuevo(&almacenFalso{})

	_, err := servicio.Crear(context.Background(), "   ", "")
	if !errors.Is(err, ErrNombreRequerido) {
		t.Fatalf("se esperaba ErrNombreRequerido, se obtuvo %v", err)
	}
}

func TestCrearRechazaNombreMuyLargo(t *testing.T) {
	servicio := Nuevo(&almacenFalso{})

	_, err := servicio.Crear(context.Background(), strings.Repeat("a", 121), "")
	if !errors.Is(err, ErrNombreMuyLargo) {
		t.Fatalf("se esperaba ErrNombreMuyLargo, se obtuvo %v", err)
	}
}

func TestCrearRechazaDescripcionMuyLarga(t *testing.T) {
	servicio := Nuevo(&almacenFalso{})

	_, err := servicio.Crear(context.Background(), "Nombre", strings.Repeat("a", 501))
	if !errors.Is(err, ErrDescripcionMuyLarga) {
		t.Fatalf("se esperaba ErrDescripcionMuyLarga, se obtuvo %v", err)
	}
}

func TestCrearPropagaElErrorDelAlmacen(t *testing.T) {
	servicio := Nuevo(&almacenFalso{err: errors.New("sin conexión")})

	_, err := servicio.Crear(context.Background(), "Nombre", "")
	if err == nil {
		t.Fatal("se esperaba un error")
	}
	if !strings.Contains(err.Error(), "crear ejemplo") {
		t.Fatalf("error sin contexto: %v", err)
	}
}
