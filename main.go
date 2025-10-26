package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type GetAccountInformationResponse struct {
	AccountId      int    `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountId int `json:"account_id" validate:"required,gt=0"`
}

var validate *validator.Validate

func main() {
	validate = validator.New()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&HealthCheckResponse{Status: "UP"})
	})
	r.Get("/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if accountId := chi.URLParam(r, "accountId"); accountId == "" {
			render.Render(w, r, ErrNotFound)
			return
		}

		// TODO: Implement

		json.NewEncoder(w).Encode(&GetAccountInformationResponse{
			AccountId:      123,
			DocumentNumber: "123",
		})
	})

	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		var input CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		if err := validate.Struct(&input); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// TODO: Implement

		w.WriteHeader(http.StatusCreated)
	})

	r.Post("/transactions", func(w http.ResponseWriter, r *http.Request) {
		var input CreateTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		if err := validate.Struct(&input); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
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

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
