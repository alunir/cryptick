package fills

import "time"

type Fill struct {
	Future    string `json:"future"`
	Market    string `json:"market"`
	Type      string `json:"type"`
	Liquidity string `json:"liquidity"`

	// only rest follow 2factor
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`

	Side string `json:"side"`

	Price   float64 `json:"price"`
	Size    float64 `json:"size"`
	Fee     float64 `json:"fee"`
	FeeRate float64 `json:"feeRate"`

	Time time.Time `json:"time"`

	ID      int `json:"id"`
	OrderID int `json:"orderId"`
	TradeID int `json:"tradeId"`
}
