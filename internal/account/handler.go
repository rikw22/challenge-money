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
	validate   *validator.Validate
	repository Repository
}

func NewHandler(validate *validator.Validate, repository Repository) *Handler {
	return &Handler{
		validate:   validate,
		repository: repository,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")
	if accountId == "" {
		render.Render(w, r, httperrors.ErrNotFound)
		return
	}

	account, err := h.repository.GetByID(r.Context(), accountId)
	if err != nil {
		render.Render(w, r, httperrors.ErrInternalServer(err))
		return
	}

	render.JSON(w, r, account)
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

	var account Account
	account.DocumentNumber = input.DocumentNumber

	err := h.repository.Create(r.Context(), &account)
	if err != nil {
		render.Render(w, r, httperrors.ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &GetResponse{
		AccountId:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}
