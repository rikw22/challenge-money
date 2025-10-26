package account

type CreateRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type GetResponse struct {
	AccountId      int    `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
