package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/api"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/config"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/repository"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/store"
)

var version = "dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Cargar()
	if err != nil {
		logger.Error("configuración inválida", "error", err)
		os.Exit(1)
	}

	zona, err := time.LoadLocation(cfg.ZonaHoraria)
	if err != nil {
		logger.Error("zona horaria inválida", "error", err)
		os.Exit(1)
	}
	time.Local = zona

	ctx, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelar()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("no se pudo crear la conexión a PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("PostgreSQL no responde al arrancar", "error", err)
		os.Exit(1)
	}

	consultas := repository.New(pool)
	almacen := store.Nuevo(consultas)
	servicio := service.Nuevo(almacen)

	handler, _ := api.Construir(cfg, servicio, pool, logger, version)
	servidor := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Puerto),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	contextoSenales, detener := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer detener()

	go func() {
		<-contextoSenales.Done()
		ctxCierre, cancelarCierre := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelarCierre()
		_ = servidor.Shutdown(ctxCierre)
	}()

	logger.Info("API iniciada",
		"servicio", "pingeso-grupo5",
		"version", version,
		"puerto", cfg.Puerto,
		"entorno", cfg.Entorno,
	)

	if err := servidor.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("la API se detuvo", "error", err)
		os.Exit(1)
	}
}
