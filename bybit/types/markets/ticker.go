package markets

import "time"

type Ticker struct {
	Bid     float64   `json:"bid"`
	Ask     float64   `json:"ask"`
	BidSize float64   `json:"bidSize"`
	AskSize float64   `json:"askSize"`
	Time    time.Time `json:"time"`
}
