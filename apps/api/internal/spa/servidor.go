package spa

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type archivo struct {
	datos      []byte
	modificado time.Time
}

func Nuevo(directorio, backend string) (http.Handler, error) {
	destino, err := url.Parse(backend)
	if err != nil || destino.Host == "" || (destino.Scheme != "http" && destino.Scheme != "https") {
		return nil, errors.New("BACKEND_URL inválido")
	}
	// El filesystem solo se consulta al arrancar, antes de recibir peticiones.
	raiz, err := os.OpenRoot(directorio)
	if err != nil {
		return nil, errors.New("directorio estático no disponible")
	}
	defer raiz.Close()
	archivos := make(map[string]archivo)
	err = fs.WalkDir(raiz.FS(), ".", func(nombre string, entrada fs.DirEntry, fallo error) error {
		if fallo != nil {
			return fallo
		}
		if entrada.IsDir() {
			return nil
		}
		if !entrada.Type().IsRegular() {
			return errors.New("el bundle debe contener solo archivos regulares")
		}
		datos, fallo := fs.ReadFile(raiz.FS(), nombre)
		if fallo != nil {
			return fallo
		}
		info, fallo := entrada.Info()
		if fallo != nil {
			return fallo
		}
		archivos["/"+nombre] = archivo{datos, info.ModTime()}
		return nil
	})
	if err != nil {
		return nil, errors.New("bundle estático inválido")
	}
	if _, existe := archivos["/index.html"]; !existe {
		return nil, errors.New("falta index.html")
	}
	proxy := httputil.NewSingleHostReverseProxy(destino)
	proxy.Transport = &http.Transport{ResponseHeaderTimeout: 5 * time.Second}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Backend no disponible", http.StatusBadGateway)
	}
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
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		for _, segmento := range strings.Split(r.URL.Path, "/") {
			if segmento == ".." || strings.Contains(segmento, `\`) {
				http.Error(w, "Ruta inválida", http.StatusBadRequest)
				return
			}
		}
		nombre := path.Clean("/" + r.URL.Path)
		contenido, existe := archivos[nombre]
		if strings.HasPrefix(nombre, "/assets/") {
			if !existe {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if !existe {
			nombre = "/index.html"
			contenido = archivos[nombre]
		}
		http.ServeContent(w, r, nombre, contenido.modificado, bytes.NewReader(contenido.datos))
	}), nil
}
