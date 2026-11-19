package api

import (
	"encoding/json"
	"net/http"
)

type problema struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func escribirProblem(w http.ResponseWriter, estado int, titulo string, detalle string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(estado)
	_ = json.NewEncoder(w).Encode(problema{
		Type:   "about:blank",
		Title:  titulo,
		Status: estado,
		Detail: detalle,
	})
}
