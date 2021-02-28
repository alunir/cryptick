package markets

type OrderBook struct {
	Bids           [][]float64 `json:"bids"`
	Asks           [][]float64 `json:"asks"`
	Timestamp      uint64      `json:"timestamp"`
	ChangeId       uint64      `json:"change_id"`
	InstrumentName string      `json:"instrument_name"`
}
