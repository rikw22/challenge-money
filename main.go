package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&HealthCheckResponse{Status: "UP"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
