package realtime

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alunir/cryptick/ftx/types/fills"
	"github.com/alunir/cryptick/ftx/types/markets"
	"github.com/alunir/cryptick/ftx/types/orders"
	"github.com/buger/jsonparser"
	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

const (
	FTX_TICKER    = "ticker"
	FTX_TRADES    = "trades"
	FTX_ORDERBOOK = "orderbook"
	FTX_FILLS     = "fills"
	FTX_ORDERS    = "orders"
)

const (
	UNDEFINED = iota
	ERROR
	TICKER
	TRADES
	ORDERBOOK
	ORDERS
	FILLS
)

var (
	orderBookLocals = make(map[string]*markets.OrderBookLocal)
)

type request struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

// {"op": "login", "args": {"key": "<api_key>", "sign": "<signature>", "time": 1111}}
type requestForPrivate struct {
	Op   string                 `json:"op"`
	Args map[string]interface{} `json:"args"`
}

type Response struct {
	Type   int
	Symbol string

	Ticker    markets.Ticker
	Trades    []markets.Trade
	Orderbook markets.OrderBook

	Orders orders.Order
	Fills  fills.Fill

	Results error
}

func subscribe(ctx context.Context, conn *websocket.Conn, channels, symbols []string) error {
	if symbols != nil {
		for i := range channels {
			for j := range symbols {
				if val, err := json.Marshal(request{
					Op:      "subscribe",
					Channel: channels[i],
					Market:  symbols[j],
				}); err != nil {
					return err
				} else if err := conn.Write(ctx, websocket.MessageBinary, val); err != nil {
					return err
				}
			}
		}
	} else {
		for i := range channels {
			if val, err := json.Marshal(request{
				Op:      "subscribe",
				Channel: channels[i],
			}); err != nil {
				return err
			} else if err := conn.Write(ctx, websocket.MessageBinary, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func unsubscribe(ctx context.Context, conn *websocket.Conn, channels, symbols []string) error {
	if symbols != nil {
		for i := range channels {
			for j := range symbols {
				if val, err := json.Marshal(request{
					Op:      "unsubscribe",
					Channel: channels[i],
					Market:  symbols[j],
				}); err != nil {
					return err
				} else if err := conn.Write(ctx, websocket.MessageBinary, val); err != nil {
					return err
				}
			}
		}
	} else {
		for i := range channels {
			if val, err := json.Marshal(request{
				Op:      "unsubscribe",
				Channel: channels[i],
			}); err != nil {
				return err
			} else if err := conn.Write(ctx, websocket.MessageBinary, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func ping(ctx context.Context, conn *websocket.Conn) (err error) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.Ping(ctx); err != nil {
				goto EXIT
			}
		}
	}
EXIT:
	return err
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, cfg *Configuration) error {
	if cfg.l == nil {
		cfg.l = log.New(os.Stdout, "ftx websocket", log.Llongfile)
	}

	var obl *markets.OrderBookLocal
	var ok bool

RECONNECT:
	conn, _, err := websocket.Dial(ctx, cfg.url, &websocket.DialOptions{
		CompressionMode: websocket.CompressionDisabled,
	})
	if err != nil {
		return err
	}
	conn.SetReadLimit(1 << 62)

	if err := subscribe(ctx, conn, channels, symbols); err != nil {
		return err
	}

	// ping each 15sec for exchange
	go ping(ctx, conn)

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

			typeMsg, err := jsonparser.GetString(msg, "type")
			if typeMsg == "error" {
				cfg.l.Printf("[ERROR]: error: %+v", string(msg))
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			channel, err := jsonparser.GetString(msg, "channel")
			if err != nil {
				cfg.l.Printf("[ERROR]: channel error: %+v", string(msg))
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			market, err := jsonparser.GetString(msg, "market")
			if err != nil {
				cfg.l.Printf("[ERROR]: market err: %+v", string(msg))
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			res.Symbol = market

			data, _, _, err := jsonparser.Get(msg, "data")
			if err != nil {
				if isSubscribe, _ := jsonparser.GetString(msg, "type"); isSubscribe == "subscribed" {
					cfg.l.Printf("[SUCCESS]: %s %+v", isSubscribe, string(msg))
					continue
				} else {
					err = fmt.Errorf("[ERROR]: data err: %v %s", err, string(msg))
					cfg.l.Println(err)
					res.Type = ERROR
					res.Results = err
					ch <- res
					continue
				}
			}

			switch channel {
			case FTX_TICKER:
				res.Type = TICKER
				if err := json.Unmarshal(data, &res.Ticker); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal ticker %+v", err)
					continue
				}

			case FTX_TRADES:
				res.Type = TRADES
				if err := json.Unmarshal(data, &res.Trades); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal trades %+v", err)
					continue
				}

			case FTX_ORDERBOOK:
				var obr *markets.OrderBookRaw
				if err := json.Unmarshal(data, &obr); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal orderbook %+v", err)
					continue
				}

				// MEMO: 'Action' determines the received data as a snapshot or a diff-snapshot.
				// see; https://docs.ftx.com/#orderbooks
				switch obr.Action {
				case "partial":
					obl, ok = orderBookLocals[market]
					if !ok {
						obl = markets.NewOrderBookLocal()
						orderBookLocals[market] = obl
					}
					obl.LoadSnapshot(obr)
				case "update":
					obl, ok = orderBookLocals[market]
					if !ok {
						continue
					}
					obl.Update(obr)
				}

				res.Type = ORDERBOOK
				res.Orderbook = obl.GetOrderBook()
			default:
				res.Type = UNDEFINED
				res.Results = fmt.Errorf("%v", string(msg))
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

func ConnectForPrivate(ctx context.Context, ch chan Response, channels []string, cfg *Configuration) error {
	if cfg.l == nil {
		cfg.l = log.New(os.Stdout, "ftx websocket", log.Llongfile)
	}

RECONNECT:
	conn, _, err := websocket.Dial(ctx, cfg.url, nil)
	if err != nil {
		return err
	}
	conn.SetReadLimit(1 << 62)

	// sign up
	if err := signature(ctx, conn, cfg.key, cfg.secret, cfg.subaccount); err != nil {
		return err
	}

	if err := subscribe(ctx, conn, channels, nil); err != nil {
		return err
	}

	go ping(ctx, conn)

	var eg errgroup.Group
	eg.Go(func() error {
		defer conn.Close(websocket.StatusNormalClosure, "normal closure")
		defer unsubscribe(ctx, conn, channels, nil)

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

			typeMsg, err := jsonparser.GetString(msg, "type")
			if typeMsg == "error" {
				cfg.l.Printf("[ERROR]: error: %+v", string(msg))
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			channel, err := jsonparser.GetString(msg, "channel")
			if err != nil {
				cfg.l.Printf("[ERROR]: channel error: %+v", string(msg))
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			data, _, _, err := jsonparser.Get(msg, "data")
			if err != nil {
				if isSubscribe, _ := jsonparser.GetString(msg, "type"); isSubscribe == "subscribed" {
					cfg.l.Printf("[SUCCESS]: %s %+v", isSubscribe, string(msg))
					continue
				} else {
					err = fmt.Errorf("[ERROR]: data err: %v %s", err, string(msg))
					cfg.l.Println(err)
					res.Type = ERROR
					res.Results = err
					ch <- res
					continue
				}
			}

			// Private channel has not market name.
			switch channel {
			case FTX_ORDERS:
				res.Type = ORDERS
				if err := json.Unmarshal(data, &res.Orders); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal orders %+v", err)
					continue
				}

			case FTX_FILLS:
				res.Type = FILLS
				if err := json.Unmarshal(data, &res.Fills); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal fills %+v", err)
					continue
				}

			default:
				res.Type = UNDEFINED
				res.Results = fmt.Errorf("%v", string(msg))
			}

			ch <- res
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("%v", err)
	}

	goto RECONNECT
}

func signature(ctx context.Context, conn *websocket.Conn, key, secret string, subaccount []string) error {
	// key: your API key
	// time: integer current timestamp (in milliseconds)
	// sign: SHA256 HMAC of the following string, using your API secret: <time>websocket_login
	// subaccount: (optional) subaccount name
	// As an example, if:

	// time: 1557246346499
	// secret: 'Y2QTHI23f23f23jfjas23f23To0RfUwX3H42fvN-'
	// sign would be d10b5a67a1a941ae9463a60b285ae845cdeac1b11edc7da9977bef0228b96de9

	// One websocket connection may be logged in to at most one user. If the connection is already authenticated, further attempts to log in will result in 400s.

	if key == "" {
		log.Fatal("Key should be specified")
	}
	if secret == "" {
		log.Fatal("SecretKey should be specified")
	}

	msec := time.Now().UTC().UnixNano() / int64(time.Millisecond)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%dwebsocket_login", msec)))
	args := map[string]interface{}{
		"key":  key,
		"sign": hex.EncodeToString(mac.Sum(nil)),
		"time": msec,
	}
	if len(subaccount) > 0 {
		args["subaccount"] = subaccount[0]
	}

	if val, err := json.Marshal(args); err != nil {
		return err
	} else if err := conn.Write(ctx, websocket.MessageText, val); err != nil {
		return err
	}

	return nil
}
