package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/bybit/realtime"
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"orderBookL2_25"}, []string{"ETHUSD"}, realtime.Config())

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Ticker)
			case realtime.TRADES:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Trades)
			case realtime.ORDERBOOK:
				// fmt.Printf("%s	%+v\n", v.Symbol, v.Orderbook)
			case realtime.UNDEFINED:
				fmt.Printf("%s	%s\n", v.Symbol, v.Results.Error())
			}
		}
	}
}

// TODO: TestConnectPrivate after implemented
