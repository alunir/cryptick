package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/deribit/realtime"
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"quote", "trades"}, []string{"BTC-PERPETUAL"}, realtime.Config())

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
// 	go realtime.ConnectForPrivate(ctx, ch, "", "", []string{"child_order_events", "parent_order_events"}, nil)

// 	for {
// 		select {
// 		case v := <-ch:
// 			switch v.Types {
// 			case realtime.CHILD_ORDERS:
// 				fmt.Printf("%d	%+v\n", v.Types, v.ChildOrderEvent)
// 			case realtime.PARENT_ORDERS:
// 				fmt.Printf("%d	%+v\n", v.Types, v.ParentOrderEvent)
// 			case realtime.UNDEFINED:
// 				fmt.Printf("UNDEFINED %s	%s\n", v.ProductCode, v.Results.Error())
// 			}
// 		}
// 	}
// }
