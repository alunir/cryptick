package markets

import (
	"time"

	"github.com/alunir/cryptick/ftx/types"
)

// Received data types from exchange
type OrderBookRaw struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
	// Action return update/partial
	Action   string        `json:"action"`
	Time     types.FtxTime `json:"time"`
	Checksum int           `json:"checksum"`
}

// Convert data internally
type OrderBookL2 struct {
	Price float64
	Side  string
	Size  float64
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
