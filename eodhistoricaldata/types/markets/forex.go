package markets

type Forex struct {
	Symbol                string  `json:"s"`
	AskPrice              float64 `json:"a"`
	BidPrice              float64 `json:"b"`
	DailyChangePercentage string  `json:"dc"`
	DailyDifferencePrice  string  `json:"dd"`
	PrePostMarketStatus   bool    `json:"ppms"` // always false
	Timestamp             int64   `json:"t"`
}
