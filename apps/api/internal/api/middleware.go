package api

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recuperar convierte un pánico en un 500 RFC 9457 sin traza para el cliente.
// La traza solo va al log, junto al X-Request-Id que asigna el registro de cmd/api.
// http.ErrAbortHandler se relanza: es la señal estándar para abortar la respuesta.
func Recuperar(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recuperado := recover()
			if recuperado == nil {
				return
			}
			if err, ok := recuperado.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(recuperado)
			}
			slog.Error("pánico controlado",
				"request_id", w.Header().Get("X-Request-Id"),
				"ruta", r.URL.Path,
				"traza", string(debug.Stack()),
			)
			escribirProblem(w, http.StatusInternalServerError, "error interno", "ocurrió un error inesperado")
		}()
		siguiente.ServeHTTP(w, r)
	})
}
