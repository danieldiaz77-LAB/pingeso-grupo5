package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// Algoritmos de firma aceptados para el JWT del servicio de identidad (DAS §6).
var algoritmosPermitidos = []string{"EdDSA", "RS256"}

// rutasPublicas no exigen sesión: la salud la consultan Compose y los
// orquestadores sin credenciales (DAS §8).
var rutasPublicas = map[string]bool{"/healthz": true, "/readyz": true}

// Proteger exige un JWT válido en todas las rutas salvo las de salud. Sin
// jwksURL (solo dev/test, lo valida config.Cargar) devuelve el handler tal cual.
// El JWKS se refresca en segundo plano hasta que ctx termina.
func Proteger(ctx context.Context, jwksURL string, siguiente http.Handler) (http.Handler, error) {
	if strings.TrimSpace(jwksURL) == "" {
		return siguiente, nil
	}
	llaves, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, err
	}
	return conAutenticacion(llaves.Keyfunc, siguiente), nil
}

func conAutenticacion(llaves jwt.Keyfunc, siguiente http.Handler) http.Handler {
	analizador := jwt.NewParser(jwt.WithValidMethods(algoritmosPermitidos), jwt.WithExpirationRequired())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rutasPublicas[r.URL.Path] {
			siguiente.ServeHTTP(w, r)
			return
		}
		token := tokenDeSolicitud(r)
		if token == "" {
			escribirProblem(w, http.StatusUnauthorized, "no autenticado", "falta el token de sesión")
			return
		}
		if _, err := analizador.ParseWithClaims(token, &jwt.RegisteredClaims{}, llaves); err != nil {
			escribirProblem(w, http.StatusUnauthorized, "token inválido", "no se pudo validar el token")
			return
		}
		siguiente.ServeHTTP(w, r)
	})
}

func tokenDeSolicitud(r *http.Request) string {
	if cookie, err := r.Cookie("sesion"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	autorizacion := r.Header.Get("Authorization")
	if len(autorizacion) > len("bearer ") && strings.EqualFold(autorizacion[:len("bearer ")], "bearer ") {
		return strings.TrimSpace(autorizacion[len("bearer "):])
	}
	return ""
}
