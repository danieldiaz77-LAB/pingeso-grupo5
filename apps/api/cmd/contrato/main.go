package main

import (
	"flag"
	"log"
	"os"

	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/api"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/config"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
)

func main() {
	salida := flag.String("o", "", "archivo de salida; si se omite, escribe en stdout")
	flag.Parse()

	cfg := config.Config{
		Entorno:     "dev",
		FrontendURL: "http://localhost:5170",
	}
	apiHuma, _ := api.NuevaHuma(cfg, service.Nuevo(nil))

	contenido, err := apiHuma.OpenAPI().YAML()
	if err != nil {
		log.Fatalf("no se pudo generar el OpenAPI: %v", err)
	}

	if *salida == "" {
		if _, err := os.Stdout.Write(contenido); err != nil {
			log.Fatalf("no se pudo escribir el OpenAPI: %v", err)
		}
		return
	}

	if err := os.WriteFile(*salida, contenido, 0o644); err != nil {
		log.Fatalf("no se pudo escribir %s: %v", *salida, err)
	}
}
