package markets

type UsQuote struct {
	Symbol    string  `json:"s"`
	AskPrice  float64 `json:"ap"`
	BidPrice  float64 `json:"bp"`
	AskSize   float64 `json:"as"`
	BidSize   float64 `json:"bs"`
	Timestamp int64   `json:"t"`
}
