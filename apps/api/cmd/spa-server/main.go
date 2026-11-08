package main

import (
	"context"
	"errors"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/spa"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if ejecutar() != nil {
		os.Exit(1)
	}
}
func ejecutar() error {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		cliente := http.Client{Timeout: 5 * time.Second}
		r, err := cliente.Get("http://127.0.0.1:8080/healthz")
		if err != nil {
			return err
		}
		defer r.Body.Close()
		if r.StatusCode != 200 {
			return errors.New("web no saludable")
		}
		return nil
	}
	directorio := os.Getenv("STATIC_DIR")
	if directorio == "" {
		directorio = "/srv"
	}
	handler, err := spa.Nuevo(directorio, os.Getenv("BACKEND_URL"))
	if err != nil {
		return err
	}
	servidor := &http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, cancelar := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelar()
	resultado := make(chan error, 1)
	go func() { resultado <- servidor.ListenAndServe() }()
	select {
	case err := <-resultado:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		cierre, fin := context.WithTimeout(context.Background(), 10*time.Second)
		defer fin()
		return servidor.Shutdown(cierre)
	}
	return nil
}
