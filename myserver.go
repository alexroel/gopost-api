package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"name":   "Éxito",
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func datos(w http.ResponseWriter, r *http.Request) {
	// name := r.PathValue("name")
	// age := r.PathValue("age")
	// Obtener query parameters
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	data := map[string]string{
		"name":   name,
		"age":    age,
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/datos/{name}/{age}", datos)

	fmt.Println("Servidor iniciando en http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
