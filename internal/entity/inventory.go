package entity

import "time"

type InventoryItem struct {
	ID       int       `json:"id" db:"id"`
	UserID   int       `json:"userId" db:"user_id"`
	MerchID  int       `json:"merchId" db:"merch_id"`
	Quantity int       `json:"quantity" db:"quantity"`
	Date     time.Time `json:"date" db:"date"`
}
