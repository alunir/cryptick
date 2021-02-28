package types

import (
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	BUY    = "buy"
	SELL   = "sell"
	MARKET = "market"
	LIMIT  = "limit"
)

type ProductCode string

// Parse Bitflyer's time
type BitflyerTime struct {
	time.Time
}

const bitflyerTimeLayout = "2006-01-02T15:04:05.999"

// changes bitflyerTime to time.Time
func (p *BitflyerTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "Z\"")
	p.Time, err = time.Parse(bitflyerTimeLayout, s)
	return err
}
