package transaction

import (
	"net/http"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateRequest
	if err := render.DecodeJSON(r.Body, &input); err != nil {
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
