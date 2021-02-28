package markets

type Book struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type Orderbook struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []Book  `json:"bids"`
	Asks     []Book  `json:"asks"`
}
