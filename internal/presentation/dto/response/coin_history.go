package response

type CoinHistory struct {
	Received []*ReceivedTransaction `json:"received"`
	Sent     []*SentTransaction     `json:"sent"`
}
