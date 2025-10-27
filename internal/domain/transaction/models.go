package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTransactionRequest struct {
	AccountId       int     `json:"account_id" validate:"required,gt=0"`
	OperationTypeId int     `json:"operation_type_id" validate:"required,min=1,max=4"`
	Amount          float64 `json:"amount" validate:"required,gt=0,max2decimals"`
}

type CreateTransactionResponse struct {
	ID              uuid.UUID `json:"id"`
	AccountId       int       `json:"account_id"`
	OperationTypeId int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
}

type Transaction struct {
	ID              pgtype.UUID
	AccountId       int
	OperationTypeId int
	Amount          int
	EventDate       time.Time
}
