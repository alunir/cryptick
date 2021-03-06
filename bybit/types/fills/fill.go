package fills

import "time"

type Execution struct {
	Symbol      string    `json:"symbol"`
	Side        string    `json:"side"`
	OrderID     string    `json:"order_id"`
	ExecID      string    `json:"exec_id"`
	OrderLinkID string    `json:"order_link_id"`
	Price       float64   `json:"price,string"`
	OrderQty    float64   `json:"order_qty"`
	ExecType    string    `json:"exec_type"`
	ExecQty     float64   `json:"exec_qty"`
	ExecFee     float64   `json:"exec_fee,string"`
	LeavesQty   float64   `json:"leaves_qty"`
	IsMaker     bool      `json:"is_maker"`
	TradeTime   time.Time `json:"trade_time"`
}

type Position struct {
	UserID                     int64   `json:"user_id"`
	Symbol                     string  `json:"symbol"`
	Size                       float64 `json:"size"`
	Side                       string  `json:"side"`
	PositionValue              float64 `json:"position_value,string"`
	EntryPrice                 float64 `json:"entry_price,string"`
	LiqPrice                   float64 `json:"liq_price,string"`
	BustPrice                  float64 `json:"bust_price,string"`
	Leverage                   float64 `json:"leverage,string"`
	OrderMargin                float64 `json:"order_margin,string"`
	PositionMargin             float64 `json:"position_margin,string"`
	AvailableBalance           float64 `json:"available_balance,string"`
	TakeProfit                 float64 `json:"take_profit,string"`
	TakeProfitTriggerPriceType string  `json:"tp_trigger_by,string"`
	StopLoss                   float64 `json:"stop_loss,string"`
	StopLossTriggerPriceType   string  `json:"sl_trigger_by,string"`
	RealisedPnl                float64 `json:"realised_pnl,string"`
	TrailingStop               float64 `json:"trailing_stop,string"`
	TrailingActive             float64 `json:"trailing_active,string"`
	WalletBalance              float64 `json:"wallet_balance,string"`
	RiskID                     int     `json:"risk_id"`
	OccClosingFee              float64 `json:"occ_closing_fee,string"`
	OccFundingFee              float64 `json:"occ_funding_fee,string"`
	AutoAddMargin              int     `json:"auto_add_margin"`
	CumRealisedPnl             float64 `json:"cum_realised_pnl,string"`
	PositionStatus             string  `json:"position_status"`
	PositionSeq                int64   `json:"position_seq"`
}
