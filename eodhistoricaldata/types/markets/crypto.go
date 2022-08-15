package markets

type Crypto struct {
	Ticker                string `json:"s"`
	LastPrice             string `json:"p"`
	TradeSize             string `json:"q"`
	DailyChangePercentage string `json:"dc"`
	DailyDifferencePrice  string `json:"dd"`
	Timestamp             int64  `json:"t"`
}
