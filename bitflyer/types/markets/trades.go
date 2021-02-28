package markets

import (
	"github.com/alunir/cryptick/bitflyer/types"
)

type Trade struct {
	ID                         int                `json:"id"`
	Side                       string             `json:"side"`
	Price                      float64            `json:"price"`
	Size                       float64            `json:"size"`
	ExecDate                   types.BitflyerTime `json:"exec_date"`
	BuyChildOrderAcceptanceID  string             `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string             `json:"sell_child_order_acceptance_id"`
}
