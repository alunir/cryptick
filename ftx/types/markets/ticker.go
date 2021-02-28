package markets

import "github.com/alunir/cryptick/ftx/types"

type Ticker struct {
	Bid     float64       `json:"bid"`
	Ask     float64       `json:"ask"`
	BidSize float64       `json:"bidSize"`
	AskSize float64       `json:"askSize"`
	Last    float64       `json:"last"`
	Time    types.FtxTime `json:"time"`
}
