package markets

import "sort"

type Book struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type Orderbook struct {
	MidPrice float64 `json:"mid_price"`
	Bids     []Book  `json:"bids"`
	Asks     []Book  `json:"asks"`
}

func (ob *Orderbook) Sort() {
	sort.Slice(ob.Bids, func(i, j int) bool {
		return ob.Bids[i].Price > ob.Bids[j].Price
	})

	sort.Slice(ob.Asks, func(i, j int) bool {
		return ob.Asks[i].Price < ob.Asks[j].Price
	})

	return
}
