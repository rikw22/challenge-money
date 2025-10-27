package transaction

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateTransactionRequest
	if err := render.DecodeJSON(r.Body, &input); err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	var t Transaction
	t.AccountId = input.AccountId
	t.OperationTypeId = input.OperationTypeId
	t.Amount = int(input.Amount * 100)
	t.EventDate = time.Now()

	err := h.repository.Create(r.Context(), &t)
	if err != nil {
		render.Render(w, r, httperrors.ErrInternalServer(err))
		return
	}

	responseID, err := uuid.FromBytes(t.ID.Bytes[:])
	if err != nil {
		render.Render(w, r, httperrors.ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &CreateTransactionResponse{
		ID:              responseID,
		AccountId:       t.AccountId,
		OperationTypeId: t.OperationTypeId,
		Amount:          float64(t.Amount) / 100,
	})
}
