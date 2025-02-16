package request

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required,min=4,max=30"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}
