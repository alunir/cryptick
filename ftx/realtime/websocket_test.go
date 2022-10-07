package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/ftx/realtime"
)

var (
	cfg = realtime.Config(realtime.Key(""), realtime.SecretKey(""))
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{realtime.FTX_TRADES}, []string{"ETH-PERP"}, cfg)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Ticker)
			case realtime.TRADES:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Trades)
				for i := range v.Trades {
					if v.Trades[i].Liquidation {
						fmt.Printf("-----------------------------%+v\n", v.Trades[i])
					}
				}
			case realtime.ORDERBOOK:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Orderbook)

			case realtime.UNDEFINED:
				fmt.Printf("%s	%s\n", v.Symbol, v.Results.Error())
			}
		}
	}

}

func TestConnectForPrivate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.ConnectForPrivate(ctx, ch, []string{"orders", "fills"}, cfg)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.ORDERS:
				fmt.Printf("%d	%+v\n", v.Type, v.Orders)
			case realtime.FILLS:
				fmt.Printf("%d	%+v\n", v.Type, v.Fills)

			case realtime.UNDEFINED:
				fmt.Printf("UNDEFINED %s	%s\n", v.Symbol, v.Results.Error())
			}
		}
	}
}
