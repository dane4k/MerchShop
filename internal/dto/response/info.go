package response

import "MerchShop/internal/dto"

type InfoResponse struct {
	Coins       int                  `json:"coins"`
	Inventory   []*dto.InventoryItem `json:"inventory"`
	CoinHistory *dto.CoinHistory     `json:"coinHistory"`
}
