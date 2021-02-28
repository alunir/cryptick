package orders

import "time"

type Order struct {
	ID            int       `json:"id"`
	ClientID      string    `json:"clientId"`
	Market        string    `json:"market"`
	Type          string    `json:"type"`
	Side          string    `json:"side"`
	Size          float64   `json:"size"`
	Price         float64   `json:"price"`
	ReduceOnly    bool      `json:"reduceOnly"`
	Ioc           bool      `json:"ioc"`
	PostOnly      bool      `json:"postOnly"`
	Status        string    `json:"status"`
	FilledSize    float64   `json:"filledSize"`
	RemainingSize float64   `json:"remainingSize"`
	AvgFillPrice  float64   `json:"avgFillPrice"`
	CreatedAt     time.Time `json:"createdAt"`
}
