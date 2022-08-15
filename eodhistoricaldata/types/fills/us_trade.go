package fills

type UsTrade struct {
	Ticker                      string  `json:"s"`
	Price                       float64 `json:"p"`
	DarkPool                    bool    `json:"dp"`
	Timestamp                   int64   `json:"t"`
	CandleConstructionParameter string  `json:"c"`
}
