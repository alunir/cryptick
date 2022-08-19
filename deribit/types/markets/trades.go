package markets

type Trade struct {
	Timestamp     uint64  `json:"timestamp"`
	Symbol        string  `json:"instrument_name"`
	Side          string  `json:"direction,omitempty"`
	Size          float64 `json:"amount,omitempty"`
	Price         float64 `json:"price,omitempty"`
	TickDirection uint8   `json:"tick_direction,omitempty"`
	TradeID       string  `json:"trade_id,omitempty"`
	TradeSeq      uint64  `json:"trade_seq,omitempty"`
	MarkPrice     float64 `json:"mark_price,omitempty"`
	IndexPrice    float64 `json:"index_price,omitempty"`
}
