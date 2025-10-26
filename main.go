package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rikw22/challenge-money/internal/account"
	"github.com/rikw22/challenge-money/internal/database"
	"github.com/rikw22/challenge-money/internal/health"
	"github.com/rikw22/challenge-money/internal/transaction"
)

var validate *validator.Validate

func main() {
	// DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:65432/postgres?sslmode=disable"
		log.Println("DATABASE_URL environment variable is not set, failing back to ", dbURL)
	}

	dbPool, err := database.NewConnection(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbPool.Close()

	accountRepo := account.NewRepository(dbPool)

	// HTTP
	validate = validator.New()

	healthHandler := health.NewHandler()
	accountHandler := account.NewHandler(validate, accountRepo)
	transactionHandler := transaction.NewHandler(validate)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/health", healthHandler.Check)
	r.Post("/accounts", accountHandler.Create)
	r.Get("/accounts/{accountId}", accountHandler.Get)
	r.Post("/transactions", transactionHandler.Create)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
