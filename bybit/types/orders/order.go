package orders

import "time"

type Order struct {
	OrderID        string    `json:"order_id"`
	OrderLinkID    string    `json:"order_link_id"`
	Symbol         string    `json:"symbol"`
	Side           string    `json:"side"`
	OrderType      string    `json:"order_type"`
	Price          float64   `json:"price,string"`
	Qty            float64   `json:"qty"`
	TimeInForce    string    `json:"time_in_force"` // GoodTillCancel/ImmediateOrCancel/FillOrKill/PostOnly
	CreateType     string    `json:"create_type"`
	CancelType     string    `json:"cancel_type"`
	OrderStatus    string    `json:"order_status"`
	LeavesQty      float64   `json:"leaves_qty"`
	CumExecQty     float64   `json:"cum_exec_qty"`
	CumExecValue   float64   `json:"cum_exec_value,string"`
	CumExecFee     float64   `json:"cum_exec_fee,string"`
	Timestamp      time.Time `json:"timestamp"`
	TakeProfit     float64   `json:"take_profit,string"`
	StopLoss       float64   `json:"stop_loss,string"`
	TrailingStop   float64   `json:"trailing_stop,string"`
	TrailingActive float64   `json:"trailing_active,string"`
	LastExecPrice  float64   `json:"last_exec_price,string"`
	ReduceOnly     bool      `json:"reduce_only,bool"`
	CloseOnTrigger bool      `json:"close_on_trigger,bool"`
}
