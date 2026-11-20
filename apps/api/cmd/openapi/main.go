package main

import (
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/api"
	"github.com/danieldiaz77-LAB/pingeso-grupo5/apps/api/internal/service"
	"os"
)

func main() {
	_, contrato := api.Nueva(nil, "dev", service.Nuevo(nil))
	datos, err := contrato.OpenAPI().YAML()
	if err != nil {
		panic(err)
	}
	if _, err = os.Stdout.Write(datos); err != nil {
		panic(err)
	}
}
