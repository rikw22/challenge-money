package account

import "time"

type CreateRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type CreateResponse struct {
	ID             int    `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type GetResponse struct {
	ID             int    `json:"account_id"`
	DocumentNumber string `json:"document_number"`
	CreatedAt      string `json:"created_at"`
}

type Account struct {
	ID             int
	DocumentNumber string
	CreatedAt      time.Time
}
