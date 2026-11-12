package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func registrarSalud(mux *http.ServeMux, salud VerificadorBD) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		escribirJSON(w, http.StatusOK, map[string]string{
			"estado": "vivo",
		})
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancelar := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancelar()

		if err := salud.Ping(ctx); err != nil {
			escribirJSON(w, http.StatusServiceUnavailable, map[string]string{
				"estado":  "no listo",
				"detalle": "PostgreSQL no responde",
			})
			return
		}

		escribirJSON(w, http.StatusOK, map[string]string{
			"estado":  "listo",
			"detalle": "PostgreSQL responde",
		})
	})
}

func escribirJSON(w http.ResponseWriter, estado int, valor any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(estado)
	_ = json.NewEncoder(w).Encode(valor)
}
