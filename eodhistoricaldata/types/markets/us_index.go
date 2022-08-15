package markets

type UsIndex struct {
	Ticker                string  `json:"s"`
	LastPrice             float64 `json:"p"`
	DailyChangePercentage string  `json:"dc"`
	DailyDifferencePrice  string  `json:"dd"`
	PrePostMarketStatus   string  `json:"ppms"`
	Timestamp             int64   `json:"t"`
}
