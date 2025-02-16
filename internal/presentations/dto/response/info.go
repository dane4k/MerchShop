package response

type InfoResponse struct {
	Coins       int              `json:"coins"`
	Inventory   []*InventoryItem `json:"inventory"`
	CoinHistory *CoinHistory     `json:"coinHistory"`
}
