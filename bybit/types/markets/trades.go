package markets

import "time"

type Trade struct {
	Time          time.Time `json:"timestamp"`
	TradeTimeMs   string    `json:"trade_time_ms"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Size          float64   `json:"size"`
	Price         string    `json:"price"`
	TickDirection string    `json:"tick_direction"`
	TradeID       string    `json:"trade_id"`
	CrossSeq      int       `json:"cross_seq"`
	IsBlockTrade  string    `json:"is_block_trade"`
}
