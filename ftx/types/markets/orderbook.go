package markets

import "github.com/alunir/cryptick/ftx/types"

type Orderbook struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
	// Action return update/partial
	Action   string        `json:"action"`
	Time     types.FtxTime `json:"time"`
	Checksum int           `json:"checksum"`
}
