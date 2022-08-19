package markets

type UsTrade struct {
	Symbol         string           `json:"s"`
	Price          float64          `json:"p"`
	Size           float64          `json:"v"`
	DarkPool       bool             `json:"dp"`
	Timestamp      int64            `json:"t"`
	TradeCondition []TradeCondition `json:"c"`
}

type TradeCondition uint8

// TradeCondition represents the consitions of the transaction.
// https://polygon.io/glossary/us/stocks/trade-conditions
const (
	REGULAR_SALE                  TradeCondition = 0
	ACQUISITION                   TradeCondition = 1
	AVERAGE_PRICE_TRADE           TradeCondition = 2
	AUTOMATIC_EXECUTION           TradeCondition = 3
	BUNCHED_TRADE                 TradeCondition = 4
	BUNCHED_SOLD_TRADE            TradeCondition = 5
	CASH_SALE                     TradeCondition = 7
	CLOSING_POINTS                TradeCondition = 8
	CROSS_TRADE                   TradeCondition = 9
	DERIVATIVELY_PRICED           TradeCondition = 10
	DISTRIBUTION                  TradeCondition = 11
	FORM_T_TRADE                  TradeCondition = 12
	EXTENDED_TRADING_HOURS        TradeCondition = 13
	INTERMARKET_SWEEP             TradeCondition = 14
	MARKET_CENTER_OFFICIAL_CLOSE  TradeCondition = 15
	MARKET_CENTER_OFFICIAL_OPEN   TradeCondition = 16
	MARKET_CENTER_OPENING_TRADE   TradeCondition = 17
	MARKET_CENTER_REOPENING_TRADE TradeCondition = 18
	MARKET_CENTER_CLOSING_TRADE   TradeCondition = 19
	NEXT_DAY                      TradeCondition = 20
	PRICE_VARIATION_TRADE         TradeCondition = 21
	PRIOR_REFERENCE_PRICE         TradeCondition = 22
	RULE_155_TRADE                TradeCondition = 23
	RULE_127_NYSE                 TradeCondition = 24
	OPENING_PRINTS                TradeCondition = 25
	STOPPED_STOCK                 TradeCondition = 27
	RE_OPENING_PRINTS             TradeCondition = 28
	SELLER                        TradeCondition = 29
	SOLD_LAST                     TradeCondition = 30
	SOLD                          TradeCondition = 33
	SPLIT_TRADE                   TradeCondition = 34
	STOCK_OPTION_TRADE            TradeCondition = 35
	ODD_LOT_TRADE                 TradeCondition = 37
	CORRECTED_CONSOLIDATED_CLOSE  TradeCondition = 38
	HELD                          TradeCondition = 40
	TRADE_THRU_EXEMPT             TradeCondition = 41
	CONTINGENT_TRADE              TradeCondition = 46
	CONTINGENT_TRADE_SALE         TradeCondition = 52
	QUALIFIED_CONTINGENT_TRADE    TradeCondition = 53
	PLACEHOLDER_FOR_611_EXEMPT    TradeCondition = 59
	SSR_IN_EFFECT                 TradeCondition = 60
)
