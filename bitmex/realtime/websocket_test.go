package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/bitmex/realtime"
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{realtime.BitmexWSQuote}, []string{"XBTUSD"}, realtime.Config())

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				fmt.Printf("%+v\n", v.Ticker)
			case realtime.TRADES:
				fmt.Printf("%+v\n", v.Trades)
			case realtime.ORDERBOOK:
				fmt.Printf("%+v\n", v.Orderbook)
			case realtime.ORDERBOOK_L2:
				fmt.Printf("%+v\n", v.OrderbookL2)
			case realtime.UNDEFINED:
				fmt.Printf("%s\n", v.Results.Error())
			}
		}
	}

}

// TODO: not supported yet
// func TestConnectForPrivate(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	ch := make(chan realtime.Response)
// 	go realtime.ConnectForPrivate(ctx, ch, "", "", []string{"orders", "fills"}, nil)

// 	for {
// 		select {
// 		case v := <-ch:
// 			switch v.Type {
// 			case realtime.ORDERS:
// 				fmt.Printf("%d	%+v\n", v.Type, v.Orders)
// 			case realtime.FILLS:
// 				fmt.Printf("%d	%+v\n", v.Type, v.Fills)

// 			case realtime.UNDEFINED:
// 				fmt.Printf("UNDEFINED %s	%s\n", v.Symbol, v.Results.Error())
// 			}
// 		}
// 	}
// }
