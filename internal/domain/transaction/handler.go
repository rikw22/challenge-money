package transaction

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rikw22/challenge-money/internal/domain/account"
	"github.com/rikw22/challenge-money/pkg/httperrors"
	"github.com/rikw22/challenge-money/pkg/validators"
)

type Handler struct {
	validate          *validator.Validate
	repository        Repository
	accountRepository account.Repository
}

func NewHandler(validate *validator.Validate, repository Repository, accountRepository account.Repository) *Handler {
	validate.RegisterValidation("max2decimals", validators.MaxTwoDecimals)
	return &Handler{
		validate:          validate,
		repository:        repository,
		accountRepository: accountRepository,
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

	// Validate Account ID
	exists, err := h.accountRepository.Exist(r.Context(), input.AccountId)
	if err != nil {
		render.Render(w, r, httperrors.ErrInternalServer(err))
		return
	}
	if !exists {
		render.Render(w, r, httperrors.ErrInvalidRequest(fmt.Errorf("account with id %d does not exist", input.AccountId)))
		return
	}

	// TODO: Validate OperationTypeId

	// Update the balance
	if input.OperationTypeId == 4 {
		inputAmountInt := int(input.Amount * 100)
		if err := h.dischargeNegativeBalances(r.Context(), input.AccountId, inputAmountInt); err != nil {
			render.Render(w, r, httperrors.ErrInternalServer(err))
			return
		}
	}

	// Create the transaction
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

	err = h.repository.Create(r.Context(), &t)
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

func (h *Handler) dischargeNegativeBalances(ctx context.Context, accountId int, paymentAmount int) error {
	transactionsWithNegativeBalance, err := h.repository.GetTransactionsWithNegativeBalance(ctx, accountId)
	if err != nil {
		return err
	}

	totalNegativeBalance := 0
	for _, transaction := range transactionsWithNegativeBalance {
		totalNegativeBalance += transaction.Balance
	}

	if len(transactionsWithNegativeBalance) > 0 && totalNegativeBalance < 0 {
		remainingAmount := paymentAmount
		for _, transaction := range transactionsWithNegativeBalance {
			amountOwed := -transaction.Balance
			remainingAmount = remainingAmount - amountOwed

			newBalanceValue := remainingAmount
			if remainingAmount >= 0 {
				newBalanceValue = 0
			}

			h.repository.UpdateTransactionBalance(ctx, transaction.ID, newBalanceValue)

			if remainingAmount < 0 {
				break
			}
		}
	}

	return nil
}
