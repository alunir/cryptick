package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/bybit/realtime"
)

var (
	cfg      = realtime.Config(realtime.Key(""), realtime.SecretKey(""))
	cfg_usdt = realtime.Config(realtime.Key(""), realtime.SecretKey(""), realtime.EndpointOption(realtime.USDT_MAINNET1_PUBLIC))
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{realtime.BybitWSTrade}, []string{"ETHUSD", "BTCUSD"}, cfg)

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

func TestConnectUSDT(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)

	// channels should be single string in USDT MAINNET
	go realtime.Connect(ctx, ch, []string{realtime.BybitWSTrade}, []string{"BTCUSDT"}, cfg_usdt)

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

func TestConnectForPrivate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.ConnectForPrivate(ctx, ch, []string{"order", "execution", "position"}, cfg)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.ORDERS:
				fmt.Printf("ORDERS(%d) %+v\n", v.Type, v.Orders)
			case realtime.FILLS:
				fmt.Printf("FILLS(%d) %+v\n", v.Type, v.Fills)
			case realtime.POSITIONS:
				fmt.Printf("POSITIONS(%d) %+v\n", v.Type, v.Positions)

			case realtime.UNDEFINED:
				fmt.Printf("UNDEFINED(%d) %s\n", v.Type, v.Results.Error())
			}
		}
	}
}
