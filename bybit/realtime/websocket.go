package realtime

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/alunir/cryptick/bybit/types/fills"
	"github.com/alunir/cryptick/bybit/types/markets"
	"github.com/alunir/cryptick/bybit/types/orders"
	"github.com/buger/jsonparser"
	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

const (
	BybitWSOrderBookL2_25  = "orderBookL2_25"        // also publish Ticker
	BybitWSOrderBookL2_200 = "orderBookL2_25.100ms." // Not supported yet
	BybitWSKLine           = "kline"
	BybitWSTrade           = "trade"
	BybitWSInsurance       = "insurance"
	BybitWSInstrument      = "instrument"

	BybitWSPosition  = "position"
	BybitWSExecution = "execution"
	BybitWSOrder     = "order"

	BybitWSDisconnected = "disconnected"
)

const (
	UNDEFINED = iota
	ERROR
	TICKER
	TRADES
	ORDERBOOK
	ORDERS
	FILLS
	POSITIONS
)

var (
	orderBookLocals = make(map[string]*markets.OrderBookLocal)
)

type request struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type Response struct {
	Type      int
	Symbol    string
	Ticker    markets.Ticker
	Trades    []markets.Trade
	Orderbook markets.OrderBook
	Orders    []orders.Order
	Fills     []fills.Execution
	Positions []fills.Position
	Results   error
}

func subscribe(ctx context.Context, conn *websocket.Conn, channels, symbols []string) error {
	var message []interface{}
	var b *bytes.Buffer
	var s string
	for i := range channels {
		if symbols == nil {
			message = append(message, channels[i])
			continue
		}
		b = new(bytes.Buffer)
		for _, symbol := range symbols {
			b.WriteString(symbol)
			b.WriteString("|")
		}
		s = b.String()
		message = append(message, fmt.Sprintf("%v.%v", channels[i], s[:len(s)-1]))
	}
	if val, err := json.Marshal(request{
		Op:   "subscribe",
		Args: message,
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
		if symbols == nil {
			message = append(message, channels[i])
			continue
		}
		for j := range symbols {
			message = append(message, fmt.Sprintf("%v.%v", channels[i], symbols[j]))
		}
	}
	if val, err := json.Marshal(request{
		Op:   "unsubscribe",
		Args: message,
	}); err != nil {
		return err
	} else if err := conn.Write(ctx, websocket.MessageText, val); err != nil {
		return err
	}
	return nil
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, cfg *Configuration) error {
	if cfg.l == nil {
		cfg.l = log.New(os.Stdout, "bybit websocket", log.Llongfile)
	}

	var obl *markets.OrderBookLocal
	var ok bool
	var symbol string

RECONNECT:
	conn, _, err := websocket.Dial(ctx, cfg.url, nil)
	if err != nil {
		cfg.l.Fatal(err)
	}
	conn.SetReadLimit(1 << 62)

	if err := subscribe(ctx, conn, channels, symbols); err != nil {
		cfg.l.Fatal(err)
	}

	go ping(ctx, conn)

	var eg errgroup.Group
	eg.Go(func() error {
		defer conn.Close(websocket.StatusNormalClosure, "normal closure")
		defer unsubscribe(ctx, conn, channels, symbols)

		for {
			var res Response
			messageType, msg, err := conn.Read(ctx)
			if err != nil {
				cfg.l.Printf("[ERROR]: msg error: %+v", err)
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", err)
				ch <- res
				return fmt.Errorf("can't receive error: %v", err)
			}

			if messageType == websocket.MessageText {
				pong, _ := jsonparser.GetString(msg, "ret_msg")
				if pong == "pong" {
					continue
				}
				success, _ := jsonparser.GetBoolean(msg, "success")
				if success {
					cfg.l.Printf("[SUCCESS]: connection established")
					continue
				}
			}

			channel, err := jsonparser.GetString(msg, "topic")
			if err != nil {
				err = fmt.Errorf("[ERROR]: channel err: %v %s", err, string(msg))
				cfg.l.Println(err)
				res.Type = ERROR
				res.Results = err
				ch <- res
				continue
			}

			data, _, _, err := jsonparser.Get(msg, "data")
			if err != nil {
				err = fmt.Errorf("[ERROR]: data err: %v %s", err, string(msg))
				cfg.l.Println(err)
				res.Type = ERROR
				res.Results = err
				ch <- res
				continue
			}

			if strings.HasPrefix(channel, BybitWSOrderBookL2_25) {
				t, err := jsonparser.GetString(msg, "type")
				if err != nil {
					continue
				}
				symbol = channel[len(BybitWSOrderBookL2_25)+1:]
				switch t {
				case "snapshot":
					var ob []*markets.OrderBookL2
					if err := json.Unmarshal(data, &ob); err != nil {
						cfg.l.Printf("[WARN]: cant unmarshal orderbookL2 %+v", err)
						continue
					}

					obl, ok = orderBookLocals[symbol]
					if !ok {
						obl = markets.NewOrderBookLocal()
						orderBookLocals[symbol] = obl
					}
					obl.LoadSnapshot(ob)
				case "delta":
					var obd *markets.OrderBookL2Delta
					if err := json.Unmarshal(data, &obd); err != nil {
						cfg.l.Printf("[WARN]: cant unmarshal orderbookL2Delta %+v", err)
						continue
					}

					obl, ok = orderBookLocals[symbol]
					if !ok {
						continue
					}
					obl.Update(obd)
				}

				var ticker markets.Ticker
				res.Type = ORDERBOOK
				res.Symbol = symbol
				res.Orderbook, ticker = obl.GetOrderBook()

				ch <- res

				// Update Ticker even if not changed the best bid/ask
				res.Type = TICKER
				res.Symbol = symbol
				res.Ticker = ticker

			} else if strings.HasPrefix(channel, BybitWSTrade) {
				if err := json.Unmarshal(data, &res.Trades); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal trades %+v", err)
					continue
				}
				res.Type = TRADES
				res.Symbol = channel[len(BybitWSTrade)+1:]
			} else if strings.HasPrefix(channel, BybitWSKLine) {
				// Not Implemented yet
				continue
			} else if strings.HasPrefix(channel, BybitWSInsurance) {
				// Not Implemented yet
				continue
			} else if strings.HasPrefix(channel, BybitWSInstrument) {
				// Not Implemented yet
				continue
			} else {
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
		cfg.l.Printf("%v", err)
	}

	goto RECONNECT
}

func ConnectForPrivate(ctx context.Context, ch chan Response, channels []string, cfg *Configuration) {
	if cfg.l == nil {
		cfg.l = log.New(os.Stdout, "bybit websocket", log.Llongfile)
	}

RECONNECT:
	conn, _, err := websocket.Dial(ctx, cfg.url, nil)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetReadLimit(1 << 62)

	// sign up
	if err := signature(ctx, conn, cfg.key, cfg.secret); err != nil {
		log.Fatal(err)
	}

	if err := subscribe(ctx, conn, channels, nil); err != nil {
		log.Fatal(err)
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

			if pong, _ := jsonparser.GetString(msg, "ret_msg"); pong == "pong" {
				continue
			}

			if success, _ := jsonparser.GetBoolean(msg, "success"); success {
				cfg.l.Printf("[SUCCESS]: private connection established")
				continue
			}

			channel, err := jsonparser.GetString(msg, "topic")
			if err != nil {
				continue
			}

			data, _, _, err := jsonparser.Get(msg, "data")
			if err != nil {
				err = fmt.Errorf("[ERROR]: data err: %v %s", err, string(msg))
				cfg.l.Println(err)
				res.Type = ERROR
				res.Results = err
				ch <- res
				continue
			}

			if strings.HasPrefix(channel, BybitWSOrder) {
				res.Type = ORDERS
				if err := json.Unmarshal(data, &res.Orders); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal orders %+v", err)
					continue
				}

			} else if strings.HasPrefix(channel, BybitWSExecution) {
				res.Type = FILLS
				if err := json.Unmarshal(data, &res.Fills); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal fills %+v", err)
					continue
				}
			} else if strings.HasPrefix(channel, BybitWSPosition) {
				res.Type = POSITIONS
				if err := json.Unmarshal(data, &res.Positions); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal fills %+v", err)
					continue
				}
			} else {
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

func signature(ctx context.Context, conn *websocket.Conn, key, secret string) error {
	if key == "" {
		log.Fatal("Key should be specified")
	}
	if secret == "" {
		log.Fatal("SecretKey should be specified")
	}

	expires := time.Now().Unix()*1000 + 10000
	req := fmt.Sprintf("GET/realtime%d", expires)
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(req))
	signature := hex.EncodeToString(sig.Sum(nil))

	if val, err := json.Marshal(request{
		Op: "auth",
		Args: []interface{}{
			key,
			//fmt.Sprintf("%v", expires),
			expires,
			signature,
		},
	}); err != nil {
		return err
	} else if err := conn.Write(ctx, websocket.MessageText, val); err != nil {
		return err
	}

	return nil
}

func ping(ctx context.Context, conn *websocket.Conn) (err error) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.Write(ctx, websocket.MessageText, []byte(`{"op":"ping"}`)); err != nil {
				goto EXIT
			}
		}
	}
EXIT:
	return err
}
