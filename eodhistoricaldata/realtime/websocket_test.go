package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alunir/cryptick/eodhistoricaldata/realtime"
	"github.com/tk42/victolinux/env"
)

var (
	cfg_forex  = realtime.Config(realtime.EndpointGroup(realtime.ENDPOINT_FOREX), realtime.Key(env.GetString("EODHISTORICALDATA_API_KEY", "")))
	cfg_crypto = realtime.Config(realtime.EndpointGroup(realtime.ENDPOINT_CRYPTO), realtime.Key(env.GetString("EODHISTORICALDATA_API_KEY", "")))
)

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	// go realtime.Connect(ctx, ch, []string{"ETH-USD", "BTC-USD"}, cfg_crypto)
	go realtime.Connect(ctx, ch, []string{"USDJPY"}, cfg_forex)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.US_QUOTE:
				fmt.Println(v.UsQuote)
			case realtime.US_TRADE:
				fmt.Println(v.UsTrade)
			case realtime.US_INDICES:
				fmt.Println(v.UsIndex)
			case realtime.FOREX:
				fmt.Println(v.Forex)
			case realtime.CRYPTO:
				fmt.Println(v.Crypto)
			case realtime.UNDEFINED:
				fmt.Printf("%s	%s\n", v.Symbol, v.Results.Error())
			}
		}
	}
}
