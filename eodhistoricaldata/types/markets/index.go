package markets

type Index struct {
	Timestamp             int64   `json:"t"`
	Symbol                string  `json:"s"`
	Price                 float64 `json:"p"`
	DailyChangePercentage string  `json:"dc"`
	DailyDifferencePrice  string  `json:"dd"`
	PrePostMarketStatus   bool    `json:"ppms"`
}
