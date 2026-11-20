//go:build e2e

package e2e

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// servidorJWKS publica una llave Ed25519 generada en la prueba y devuelve su URL
// y una función que firma tokens con la vigencia indicada (negativa = vencido).
func servidorJWKS(t *testing.T) (string, func(time.Duration) string) {
	t.Helper()
	publica, privada, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	cuerpo, err := json.Marshal(map[string]any{"keys": []map[string]string{{
		"kty": "OKP", "crv": "Ed25519", "alg": "EdDSA", "use": "sig", "kid": "llave-e2e",
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
	firmar := func(vigencia time.Duration) string {
		token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{Subject: "usuario-e2e", ExpiresAt: jwt.NewNumericDate(time.Now().Add(vigencia))})
		token.Header["kid"] = "llave-e2e"
		firmado, err := token.SignedString(privada)
		if err != nil {
			t.Fatal(err)
		}
		return firmado
	}
	return servidor.URL + "/jwks.json", firmar
}

func consultarConToken(t *testing.T, ruta, token string, codigo int) string {
	t.Helper()
	peticion, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080"+ruta, nil)
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		peticion.Header.Set("Authorization", "Bearer "+token)
	}
	cliente := http.Client{Timeout: 7 * time.Second}
	respuesta, err := cliente.Do(peticion)
	if err != nil {
		t.Fatal(err)
	}
	defer respuesta.Body.Close()
	datos, err := io.ReadAll(respuesta.Body)
	if err != nil {
		t.Fatal(err)
	}
	if respuesta.StatusCode != codigo {
		t.Fatalf("%s: esperado %d, recibido %d: %s", ruta, codigo, respuesta.StatusCode, datos)
	}
	return string(datos)
}

func TestAutenticacionBinarioReal(t *testing.T) {
	base := baseDePrueba(t)
	jwks, firmar := servidorJWKS(t)
	cerrar := iniciar(t, "test", base, "AUTH_JWKS_URL="+jwks)
	defer cerrar()
	consultarConToken(t, "/healthz", "", 200)
	consultarConToken(t, "/readyz", "", 200)
	sinToken := consultarConToken(t, "/api/v1/ejemplos", "", 401)
	consultarConToken(t, "/api/v1/ejemplos", firmar(-time.Minute), 401)
	consultarConToken(t, "/api/v1/ejemplos", "no-es-jwt", 401)
	consultarConToken(t, "/api/v1/ejemplos", firmar(time.Hour), 200)
	if strings.Contains(sinToken, jwks) || !strings.Contains(sinToken, `"status":401`) {
		t.Fatalf("problema 401 inesperado: %s", sinToken)
	}
	informe := "# Autenticación JWT\n\n- Salud sin token: healthz 200, readyz 200\n- Ejemplos sin token: 401 problem+json\n- Token vencido o mal formado: 401\n- Token EdDSA válido del JWKS: 200\n"
	golden, err := os.ReadFile("testdata/autenticacion.golden.md")
	if err != nil {
		t.Fatal(err)
	}
	if informe != string(golden) {
		t.Fatalf("golden distinto:\n%s", informe)
	}
}
