package transaction

type CreateRequest struct {
	AccountId int `json:"account_id" validate:"required,gt=0"`
}
