package markets

type Ticker struct {
	Timestamp uint64  `json:"timestamp"`
	Symbol    string  `json:"instrument_name"`
	BidSize   float64 `json:"best_bid_amount,omitempty"`
	BidPrice  float64 `json:"best_bid_price,omitempty"`
	AskPrice  float64 `json:"best_ask_price,omitempty"`
	AskSize   float64 `json:"best_ask_amount,omitempty"`
}
