package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alunir/cryptick/deribit/types/markets"
	"github.com/buger/jsonparser"
	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

const (
	DeribitMethodSubscription = "subscription"
	DeribitMethodHeartbeat    = "heartbeat"
	DeribitChannelQuote       = "quote"
	DeribitChannelTrade       = "trades"
	DeribitChannelBook        = "book"
)

const (
	UNDEFINED = iota
	HEARTBEAT
	ERROR
	TICKER
	TRADES
	ORDERBOOK
	ORDERS
	FILLS
)

var (
	requestId int64
)

type request struct {
	Version string                 `json:"jsonrpc"`
	Id      int64                  `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type Response struct {
	Type      int
	Ticker    markets.Ticker
	Trades    []markets.Trade
	Orderbook markets.OrderBook
	// TODO: Not implemented yet
	// Orders    orders.Order
	// Fills     fills.Fill
	Results error
}

func subscribe(ctx context.Context, conn *websocket.Conn, channels, symbols []string) error {
	var message []interface{}
	for i := range channels {
		suffix := ""
		if channels[i] == DeribitChannelTrade {
			suffix = ".raw"
		} else if channels[i] == DeribitChannelBook {
			suffix = ".none.1.100ms"
		}
		for j := range symbols {
			message = append(message, fmt.Sprintf("%v.%v%v", channels[i], symbols[j], suffix))
		}
	}
	if val, err := json.Marshal(request{
		Version: "2.0",
		Method:  "public/subscribe",
		Id:      requestId,
		Params: map[string]interface{}{
			"channels": message,
		},
	}); err != nil {
		return err
	} else if err := conn.Write(ctx, websocket.MessageText, val); err != nil {
		return err
	}
	return nil
}

func unsubscribe(ctx context.Context, conn *websocket.Conn, channels, symbols []string) error {
	var message []interface{}
	for i := range channels {
		for j := range symbols {
			message = append(message, fmt.Sprintf("%v.%v", channels[i], symbols[j]))
		}
	}
	if val, err := json.Marshal(request{
		Version: "2.0",
		Method:  "public/unsubscribe",
		Id:      requestId,
		Params: map[string]interface{}{
			"channels": message,
		},
	}); err != nil {
		return err
	} else if err := conn.Write(ctx, websocket.MessageText, val); err != nil {
		return err
	}
	return nil
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, cfg *Configuration) error {
	if cfg.l == nil {
		cfg.l = log.New(os.Stdout, "deribit websocket", log.Llongfile)
	}

RECONNECT:
	conn, _, err := websocket.Dial(ctx, cfg.url, nil)
	if err != nil {
		return err
	}
	conn.SetReadLimit(1 << 62)

	requestId = 1
	if err := setHeartbeat(ctx, conn); err != nil {
		return err
	}

	if err := subscribe(ctx, conn, channels, symbols); err != nil {
		return err
	}

	var eg errgroup.Group
	eg.Go(func() error {
		defer conn.Close(websocket.StatusNormalClosure, "normal closure")
		defer unsubscribe(ctx, conn, channels, symbols)

		for {
			var res Response
			_, msg, err := conn.Read(ctx)
			if err != nil {
				cfg.l.Printf("[ERROR]: msg error: %+v", err)
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", err)
				ch <- res
				return fmt.Errorf("can't receive error: %v", err)
			}

			result, err := jsonparser.GetString(msg, "result")
			if result == "ok" {
				cfg.l.Printf("[SUCCESS]: %+v", string(msg))
				continue
			}

			method, _ := jsonparser.GetString(msg, "method")
			switch method {
			case DeribitMethodSubscription:
				channelFull, err := jsonparser.GetString(msg, "params", "channel")
				if err != nil {
					err = fmt.Errorf("[ERROR]: channel err: %v %s", err, string(msg))
					cfg.l.Println(err)
					res.Type = ERROR
					res.Results = err
					ch <- res
					continue
				}

				data, _, _, err := jsonparser.Get(msg, "params", "data")
				if err != nil {
					err = fmt.Errorf("[ERROR]: data err: %v %s", err, string(msg))
					cfg.l.Println(err)
					res.Type = ERROR
					res.Results = err
					ch <- res
					continue
				}

				// MEMO: quote.BTC-PERPETUAL -> quote
				switch channelFull[:strings.Index(channelFull, ".")] {
				case DeribitChannelQuote:
					res.Type = TICKER
					if err := json.Unmarshal(data, &res.Ticker); err != nil {
						cfg.l.Printf("[WARN]: cant unmarshal ticker %+v", err)
						continue
					}
				case DeribitChannelTrade:
					res.Type = TRADES
					if err := json.Unmarshal(data, &res.Trades); err != nil {
						cfg.l.Printf("[WARN]: cant unmarshal trades %+v", err)
						continue
					}
				case DeribitChannelBook:
					res.Type = ORDERBOOK
					if err := json.Unmarshal(data, &res.Orderbook); err != nil {
						cfg.l.Printf("[WARN]: cant unmarshal orderbook %+v", err)
						continue
					}
				default:
					res.Type = UNDEFINED
					res.Results = fmt.Errorf("%v", string(msg))
				}
			case DeribitMethodHeartbeat:
				err := testRequest(ctx, conn)
				if err != nil {
					cfg.l.Printf("[ERROR]: error: %+v", err)
					res.Type = ERROR
					res.Results = fmt.Errorf("%v", err)
					ch <- res
					return fmt.Errorf("can't send heartbeat error: %v", err)
				}
				res.Type = HEARTBEAT
			default:
				cfg.l.Printf("[NOTIFY]: %+v", string(msg))
				continue
			}

			select { // 外部からの停止
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			ch <- res
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("%v", err)
	}

	goto RECONNECT
}

func setHeartbeat(ctx context.Context, conn *websocket.Conn) error {
	if err := conn.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf(`{"jsonrpc": "2.0", "method": "public/set_heartbeat", "id": %v, "params": {"interval": 60}}`, requestId))); err != nil {
		return fmt.Errorf("Failed to send heartbeat")
	}
	return nil
}

func testRequest(ctx context.Context, conn *websocket.Conn) error {
	if err := conn.Write(ctx, websocket.MessageText, []byte(`{"method": "public/test", "params": {}}`)); err != nil {
		return fmt.Errorf("Failed to send heartbeat")
	}
	return nil
}
