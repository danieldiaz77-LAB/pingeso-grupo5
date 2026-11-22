package api

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func servidorJWKS(t *testing.T, publica ed25519.PublicKey) *httptest.Server {
	t.Helper()
	cuerpo, err := json.Marshal(map[string]any{"keys": []map[string]string{{
		"kty": "OKP", "crv": "Ed25519", "alg": "EdDSA", "use": "sig", "kid": "llave-prueba",
		"x": base64.RawURLEncoding.EncodeToString(publica),
	}}})
	if err != nil {
		t.Fatal(err)
	}
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(cuerpo)
	}))
	t.Cleanup(servidor.Close)
	return servidor
}

func firmar(t *testing.T, privada ed25519.PrivateKey, reclamos jwt.Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, reclamos)
	token.Header["kid"] = "llave-prueba"
	firmado, err := token.SignedString(privada)
	if err != nil {
		t.Fatal(err)
	}
	return firmado
}

func TestProtegerExigeTokenValidoSalvoSalud(t *testing.T) {
	publica, privada, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, otraPrivada, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	jwks := servidorJWKS(t, publica)
	ctx, cancelar := context.WithCancel(context.Background())
	defer cancelar()
	destino := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	protegido, err := Proteger(ctx, jwks.URL, destino)
	if err != nil {
		t.Fatal(err)
	}

	vigente := jwt.RegisteredClaims{Subject: "usuario", ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}
	vencido := jwt.RegisteredClaims{Subject: "usuario", ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour))}
	compartida := make([]byte, 32)
	if _, err = rand.Read(compartida); err != nil {
		t.Fatal(err)
	}
	simetrico, err := jwt.NewWithClaims(jwt.SigningMethodHS256, vigente).SignedString(compartida)
	if err != nil {
		t.Fatal(err)
	}

	for _, caso := range []struct {
		nombre, ruta, cookie, autorizacion string
		estado                             int
	}{
		{"healthz público", "/healthz", "", "", http.StatusNoContent},
		{"readyz público", "/readyz", "", "", http.StatusNoContent},
		{"sin token", "/api/v1/ejemplos", "", "", http.StatusUnauthorized},
		{"bearer vacío", "/api/v1/ejemplos", "", "Bearer ", http.StatusUnauthorized},
		{"esquema distinto", "/api/v1/ejemplos", "", "Basic abc", http.StatusUnauthorized},
		{"token mal formado", "/api/v1/ejemplos", "", "Bearer no-es-jwt", http.StatusUnauthorized},
		{"vencido", "/api/v1/ejemplos", "", "Bearer " + firmar(t, privada, vencido), http.StatusUnauthorized},
		{"sin exp", "/api/v1/ejemplos", "", "Bearer " + firmar(t, privada, jwt.RegisteredClaims{Subject: "usuario"}), http.StatusUnauthorized},
		{"otra llave", "/api/v1/ejemplos", "", "Bearer " + firmar(t, otraPrivada, vigente), http.StatusUnauthorized},
		{"algoritmo no permitido", "/api/v1/ejemplos", "", "Bearer " + simetrico, http.StatusUnauthorized},
		{"ruta de salud con sufijo", "/healthz/extra", "", "", http.StatusUnauthorized},
		{"bearer válido", "/api/v1/ejemplos", "", "bearer " + firmar(t, privada, vigente), http.StatusNoContent},
		{"cookie válida", "/api/v1/ejemplos", firmar(t, privada, vigente), "", http.StatusNoContent},
	} {
		t.Run(caso.nombre, func(t *testing.T) {
			peticion := httptest.NewRequest(http.MethodGet, caso.ruta, nil)
			if caso.cookie != "" {
				peticion.AddCookie(&http.Cookie{Name: "sesion", Value: caso.cookie})
			}
			if caso.autorizacion != "" {
				peticion.Header.Set("Authorization", caso.autorizacion)
			}
			respuesta := httptest.NewRecorder()
			protegido.ServeHTTP(respuesta, peticion)
			if respuesta.Code != caso.estado {
				t.Fatalf("esperado %d, recibido %d: %s", caso.estado, respuesta.Code, respuesta.Body)
			}
			if caso.estado == http.StatusUnauthorized {
				if respuesta.Header().Get("Content-Type") != "application/problem+json" {
					t.Fatalf("content-type: %q", respuesta.Header().Get("Content-Type"))
				}
				if strings.Contains(respuesta.Body.String(), "llave-prueba") || strings.Contains(respuesta.Body.String(), jwks.URL) {
					t.Fatalf("el error expone detalles internos: %s", respuesta.Body)
				}
			}
		})
	}
}

func TestProtegerSinJWKSNoEnvuelve(t *testing.T) {
	destino := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) })
	protegido, err := Proteger(context.Background(), "  ", destino)
	if err != nil {
		t.Fatal(err)
	}
	respuesta := httptest.NewRecorder()
	protegido.ServeHTTP(respuesta, httptest.NewRequest(http.MethodGet, "/api/v1/ejemplos", nil))
	if respuesta.Code != http.StatusTeapot {
		t.Fatalf("estado %d", respuesta.Code)
	}
}
