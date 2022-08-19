package markets

type UsTrade struct {
	Symbol                      string  `json:"s"`
	Price                       float64 `json:"p"`
	Size                        float64 `json:"v"`
	DarkPool                    bool    `json:"dp"`
	Timestamp                   int64   `json:"t"`
	CandleConstructionParameter string  `json:"c"`
}
