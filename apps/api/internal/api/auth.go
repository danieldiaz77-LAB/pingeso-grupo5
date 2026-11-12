package api

import (
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

func conAutenticacion(jwksURL string, next http.Handler) (http.Handler, error) {
	if strings.TrimSpace(jwksURL) == "" {
		return next, nil
	}

	llaves, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tokenDeSolicitud(r)
		if token == "" {
			escribirProblem(w, http.StatusUnauthorized, "no autenticado", "falta el token de sesión")
			return
		}

		reclamos := jwt.MapClaims{}
		if _, err := jwt.ParseWithClaims(token, reclamos, llaves.Keyfunc); err != nil {
			escribirProblem(w, http.StatusUnauthorized, "token inválido", "no se pudo validar el token")
			return
		}

		next.ServeHTTP(w, r)
	}), nil
}

func tokenDeSolicitud(r *http.Request) string {
	if cookie, err := r.Cookie("sesion"); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	autorizacion := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(autorizacion), "bearer ") {
		return strings.TrimSpace(autorizacion[len("bearer "):])
	}

	return ""
}
