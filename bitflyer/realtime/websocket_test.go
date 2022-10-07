package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/bitflyer/realtime"
)

var (
	cfg = realtime.Config(realtime.Key(""), realtime.SecretKey(""))
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{realtime.BITFLYER_TICKER}, []string{"FX_BTC_JPY"}, cfg)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				fmt.Printf("%s	%+v\n", v.ProductCode, v.Ticker)
			case realtime.TRADES:
				fmt.Printf("%s	%+v\n", v.ProductCode, v.Trades)
			case realtime.ORDERBOOK:
				fmt.Printf("%s	%+v\n", v.ProductCode, v.Orderbook)
			case realtime.UNDEFINED:
				fmt.Printf("%s	%s\n", v.ProductCode, v.Results.Error())
			}
		}
	}

}

func TestConnectForPrivate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.ConnectForPrivate(ctx, ch, []string{realtime.BITFLYER_CHILD_ORDER_EVENTS, realtime.BITFLYER_PARENT_ORDER_EVENTS}, cfg)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.CHILD_ORDERS:
				fmt.Printf("%d	%+v\n", v.Type, v.ChildOrderEvent)
			case realtime.PARENT_ORDERS:
				fmt.Printf("%d	%+v\n", v.Type, v.ParentOrderEvent)
			case realtime.UNDEFINED:
				fmt.Printf("UNDEFINED %s	%s\n", v.ProductCode, v.Results.Error())
			}
		}
	}
}
