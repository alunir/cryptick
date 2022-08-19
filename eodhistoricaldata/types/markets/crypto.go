package markets

type Crypto struct {
	Symbol                string `json:"s"`
	Price                 string `json:"p"`
	TradeSize             string `json:"q"`
	DailyChangePercentage string `json:"dc"`
	DailyDifferencePrice  string `json:"dd"`
	Timestamp             int64  `json:"t"`
}
