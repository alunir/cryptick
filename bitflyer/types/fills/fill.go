package fills

import "time"

type ChildOrderFill struct {
	ExecID                 int    `json:"exec_id"`
	ProductCode            string `json:"product_code"`
	ChildOrderID           string `json:"child_order_id"`
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
	ChildOrderType         string `json:"child_order_type"`

	EventDate  time.Time `json:"event_date"`
	EventType  string    `json:"event_type"`
	Side       string    `json:"side"`
	Price      int       `json:"price"`
	Size       float64   `json:"size"`
	ExpireDate string    `json:"expire_date"`

	// 新設分追記
	Reason          string  `json:"reason"`
	Commission      float64 `json:"commission"`
	SFD             float64 `json:"sfd"`
	OutstandingSize float64 `json:"outstanding_size"`
}

type ParentOrderFill struct {
	ProductCode             string    `json:"product_code"`
	ParentOrderID           string    `json:"parent_order_id"`
	ParentOrderAcceptanceID string    `json:"parent_order_acceptance_id"`
	EventDate               time.Time `json:"event_date"`
	EventType               string    `json:"event_type"`
	ParentOrderType         string    `json:"parent_order_type"`
	Reason                  string    `json:"reason"`
	ParameterIndex          int       `json:"parameter_index"`
	ChildOrderType          string    `json:"child_order_type"`
	Side                    string    `json:"side"`
	Price                   int       `json:"price"`
	Size                    float64   `json:"size"`
	ExpireDate              time.Time `json:"expire_date"`
	ChildOrderAcceptanceID  string    `json:"child_order_acceptance_id"`
}
