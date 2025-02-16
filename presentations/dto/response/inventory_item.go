package response

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}
