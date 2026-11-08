package spa

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Nuevo(directorio, backend string) (http.Handler, error) {
	destino, err := url.Parse(backend)
	if err != nil || destino.Host == "" || (destino.Scheme != "http" && destino.Scheme != "https") {
		return nil, errors.New("BACKEND_URL inválido")
	}
	proxy := httputil.NewSingleHostReverseProxy(destino)
	proxy.Transport = &http.Transport{ResponseHeaderTimeout: 5 * time.Second}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Backend no disponible", http.StatusBadGateway)
	}
	archivos := http.FileServer(http.Dir(directorio))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		if r.URL.Path == "/healthz" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"estado":"ok"}`))
			return
		}
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			proxy.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			archivos.ServeHTTP(w, r)
			return
		}
		ruta := filepath.Join(directorio, filepath.FromSlash(strings.TrimPrefix(filepath.Clean("/"+r.URL.Path), "/")))
		if info, err := os.Stat(ruta); err != nil || info.IsDir() {
			http.ServeFile(w, r, filepath.Join(directorio, "index.html"))
			return
		}
		archivos.ServeHTTP(w, r)
	}), nil
}
