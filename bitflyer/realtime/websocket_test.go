package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/bitflyer/realtime"
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"lightning_ticker"}, []string{"FX_BTC_JPY"}, nil)

	for {
		select {
		case v := <-ch:
			switch v.Types {
			case realtime.TICKER:
				fmt.Printf("%s	%+v\n", v.ProductCode, v.Ticker)

			case realtime.EXECUTIONS:
				fmt.Printf("%s	%+v\n", v.ProductCode, v.Executions)
			case realtime.BOARD:
				fmt.Printf("%s	%+v\n", v.ProductCode, v.Board)

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
	go realtime.ConnectForPrivate(ctx, ch, "", "", []string{"child_order_events", "parent_order_events"}, nil)

	for {
		select {
		case v := <-ch:
			switch v.Types {
			case realtime.CHILD_ORDERS:
				fmt.Printf("%d	%+v\n", v.Types, v.ChildOrderEvent)
			case realtime.PARENT_ORDERS:
				fmt.Printf("%d	%+v\n", v.Types, v.ParentOrderEvent)

			case realtime.UNDEFINED:
				fmt.Printf("UNDEFINED %s	%s\n", v.ProductCode, v.Results.Error())
			}
		}
	}
}
