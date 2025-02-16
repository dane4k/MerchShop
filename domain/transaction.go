package domain

type Transaction struct {
	ID         int `json:"id" db:"id"`
	Amount     int `json:"amount" db:"amount"`
	ReceiverID int `json:"receiver_id" db:"receiver_id"`
	SenderID   int `json:"sender_id" db:"sender_id"`
}
