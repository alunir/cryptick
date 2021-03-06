package markets

import (
	"strconv"
	"time"
)

// Received data types from exchange
type OrderBookL2 struct {
	ID     int64   `json:"id"`
	Price  float64 `json:"price,string"`
	Side   string  `json:"side"`
	Size   int64   `json:"size"`
	Symbol string  `json:"symbol"`
}

type OrderBookL2Delta struct {
	Delete []*OrderBookL2 `json:"delete"`
	Update []*OrderBookL2 `json:"update"`
	Insert []*OrderBookL2 `json:"insert"`
}

func (o *OrderBookL2) Key() string {
	return strconv.FormatInt(o.ID, 10)
}

// Transformed data types
type Item struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type OrderBook struct {
	Bids      []Item    `json:"bids"`
	Asks      []Item    `json:"asks"`
	Timestamp time.Time `json:"timestamp"`
}
