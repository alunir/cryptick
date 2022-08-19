package markets

type Index struct {
	Symbol                string  `json:"s"`
	Price                 float64 `json:"p"`
	DailyChangePercentage string  `json:"dc"`
	DailyDifferencePrice  string  `json:"dd"`
	PrePostMarketStatus   string  `json:"ppms"`
	Timestamp             int64   `json:"t"`
}
