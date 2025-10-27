package transaction

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rikw22/challenge-money/pkg/httperrors"
	"github.com/rikw22/challenge-money/pkg/validators"
)

type Handler struct {
	validate   *validator.Validate
	repository Repository
}

func NewHandler(validate *validator.Validate, repository Repository) *Handler {
	validate.RegisterValidation("max2decimals", validators.MaxTwoDecimals)
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

	amount := int(input.Amount * 100)
	storedAmount := amount
	// Purchases and withdrawals should be stored as negative
	if t.OperationTypeId >= 1 && t.OperationTypeId <= 3 {
		if amount > 0 {
			storedAmount = -amount
		}
	}
	t.Amount = storedAmount
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
		Amount:          float64(amount) / 100,
	})
}
