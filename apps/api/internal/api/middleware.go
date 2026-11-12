package api

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/config"
)

type claveContexto string

const claveRequestID claveContexto = "request_id"

func aplicarMiddleware(next http.Handler, cfg config.Config, logger *slog.Logger, version string) http.Handler {
	handler := conRecuperacion(next, logger, version)
	handler, err := conAutenticacion(cfg.AuthJWKSURL, handler)
	if err != nil {
		logger.Error("no se pudo configurar el JWKS", "error", err)
		handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			escribirProblem(w, http.StatusInternalServerError, "configuración inválida", "no se pudo configurar la autenticación")
		})
	}
	handler = conRegistro(handler, logger, version)
	handler = conRequestID(handler)
	return handler
}

func conRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-Id"))
		if requestID == "" {
			requestID = nuevoRequestID()
		}
		w.Header().Set("X-Request-Id", requestID)
		ctx := context.WithValue(r.Context(), claveRequestID, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func conRegistro(next http.Handler, logger *slog.Logger, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		registro := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(registro, r)

		estado := registro.status
		if estado == 0 {
			estado = http.StatusOK
		}

		logger.Info("solicitud",
			"servicio", "pingeso-grupo5",
			"version", version,
			"request_id", requestIDDesdeContext(r.Context()),
			"metodo", r.Method,
			"ruta", r.URL.Path,
			"estado", estado,
			"duracion_ms", time.Since(inicio).Milliseconds(),
		)
	})
}

func conRecuperacion(next http.Handler, logger *slog.Logger, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recuperado := recover(); recuperado != nil {
				logger.Error("pánico controlado",
					"servicio", "pingeso-grupo5",
					"version", version,
					"request_id", requestIDDesdeContext(r.Context()),
					"error", recuperado,
					"traza", string(debug.Stack()),
				)
				escribirProblem(w, http.StatusInternalServerError, "error interno", "ocurrió un error inesperado")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func requestIDDesdeContext(ctx context.Context) string {
	requestID, _ := ctx.Value(claveRequestID).(string)
	return requestID
}

func nuevoRequestID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "00000000-0000-4000-8000-000000000000"
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	)
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (r *responseWriter) WriteHeader(estado int) {
	r.status = estado
	r.ResponseWriter.WriteHeader(estado)
}

func (r *responseWriter) Write(cuerpo []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(cuerpo)
}
