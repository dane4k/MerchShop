package dto

type SentTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
