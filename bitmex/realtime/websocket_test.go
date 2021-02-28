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
	go realtime.Connect(ctx, ch, []string{realtime.BitmexWSOrderBookL2_25}, []string{"XBTUSD"}, nil)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Ticker)
			case realtime.TRADES:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Trades)
			case realtime.ORDERBOOK:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Orderbook)
			case realtime.ORDERBOOK_L2:
				fmt.Printf("%s	%+v\n", v.Symbol, v.OrderbookL2)
			case realtime.UNDEFINED:
				fmt.Printf("%s	%s\n", v.Symbol, v.Results.Error())
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
