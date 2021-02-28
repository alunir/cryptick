package markets

import "time"

type OrderBookL2 struct {
	ID     int64   `json:"id"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Size   int64   `json:"size"`
	Symbol string  `json:"symbol"`
}

// OrderBook10 contains order book 10
type OrderBook10 struct {
	Bids      [][]float64 `json:"bids"`
	Asks      [][]float64 `json:"asks"`
	Timestamp time.Time   `json:"timestamp"`
	Symbol    string      `json:"symbol"`
}
