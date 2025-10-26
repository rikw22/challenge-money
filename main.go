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

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type GetAccountInformationResponse struct {
	AccountId      int    `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountId int `json:"account_id"`
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&HealthCheckResponse{Status: "UP"})
	})
	r.Get("/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// TODO: Implement

		json.NewEncoder(w).Encode(&GetAccountInformationResponse{
			AccountId:      123,
			DocumentNumber: "123",
		})
	})

	r.Get("/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// TODO: Implement

		json.NewEncoder(w).Encode(&GetAccountInformationResponse{
			AccountId:      123,
			DocumentNumber: "123",
		})
	})

	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		var input CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// TODO: Implement

		w.WriteHeader(http.StatusCreated)
	})

	r.Post("/transactions", func(w http.ResponseWriter, r *http.Request) {
		var input CreateTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// TODO: Implement

		w.WriteHeader(http.StatusCreated)
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
