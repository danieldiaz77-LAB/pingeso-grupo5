package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/api"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/config"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/repository"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

var version = "dev"

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := ejecutar(); err != nil {
		slog.Error("arranque fallido", "causa", err.Error())
		os.Exit(1)
	}
}
func ejecutar() error {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		cliente := http.Client{Timeout: 5 * time.Second}
		r, err := cliente.Get("http://127.0.0.1:8080/healthz")
		if err != nil {
			return errors.New("API no disponible")
		}
		defer r.Body.Close()
		if r.StatusCode != 200 {
			return errors.New("API no saludable")
		}
		return nil
	}
	cfg, err := config.Cargar()
	if err != nil {
		return err
	}
	ctx, cancelar := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelar()
	base, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return errors.New("DATABASE_URL inválido")
	}
	defer base.Close()
	mux, _ := api.Nueva(base, cfg.Entorno, service.Nuevo(store.Nuevo(repository.New(base))))
	protegido, err := api.Proteger(ctx, cfg.AuthJWKSURL, mux)
	if err != nil {
		return errors.New("AUTH_JWKS_URL inválido")
	}
	servidor := &http.Server{Addr: ":8080", Handler: registrar(api.Recuperar(protegido)), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	resultado := make(chan error, 1)
	go func() { resultado <- servidor.ListenAndServe() }()
	slog.Info("API iniciada", "servicio", "base-tecnica", "version", version)
	select {
	case err := <-resultado:
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.New("no se pudo servir HTTP")
		}
	case <-ctx.Done():
		cierre, fin := context.WithTimeout(context.Background(), 10*time.Second)
		defer fin()
		return servidor.Shutdown(cierre)
	}
	return nil
}

type respuesta struct {
	http.ResponseWriter
	estado int
}

func (r *respuesta) WriteHeader(estado int) { r.estado = estado; r.ResponseWriter.WriteHeader(estado) }
func registrar(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		id := r.Header.Get("X-Request-Id")
		if id == "" || len(id) > 128 {
			bytes := make([]byte, 16)
			if _, err := rand.Read(bytes); err != nil {
				http.Error(w, "Error interno", 500)
				return
			}
			id = hex.EncodeToString(bytes)
		}
		w.Header().Set("X-Request-Id", id)
		salida := &respuesta{ResponseWriter: w, estado: 200}
		siguiente.ServeHTTP(salida, r)
		slog.Info("request", "servicio", "base-tecnica", "version", version, "request_id", id, "metodo", r.Method, "ruta", r.URL.Path, "estado", salida.estado, "duracion_ms", time.Since(inicio).Milliseconds())
	})
}
