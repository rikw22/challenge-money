package account

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rikw22/challenge-money/internal/httperrors"
)

type Handler struct {
	validate *validator.Validate
}

func NewHandler(validate *validator.Validate) *Handler {
	return &Handler{
		validate: validate,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if accountId := chi.URLParam(r, "accountId"); accountId == "" {
		render.Render(w, r, httperrors.ErrNotFound)
		return
	}

	// TODO: Implement

	json.NewEncoder(w).Encode(&GetResponse{
		AccountId:      123,
		DocumentNumber: "123",
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	// TODO: Implement

	w.WriteHeader(http.StatusCreated)
}
