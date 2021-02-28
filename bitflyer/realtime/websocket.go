package realtime

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/alunir/cryptick/bitflyer/types"
	"github.com/alunir/cryptick/bitflyer/types/fills"
	"github.com/alunir/cryptick/bitflyer/types/markets"
	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const (
	USE1                       = "wss://ws.lightstream.bitflyer.com/json-rpc"
	READDEADLINE time.Duration = 300
)

type Types int

const (
	ALL Types = iota
	TICKER
	EXECUTIONS
	BOARD
	DIFF_BOARDS
	CHILD_ORDERS
	PARENT_ORDERS
	UNDEFINED
	ERROR
)

type request struct {
	Jsonrpc string                 `json:"jsonrpc,omitempty"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id,omitempty"`
}

// {"op": "login", "args": {"key": "<api_key>", "sign": "<signature>", "time": 1111}}
type requestForPrivate struct {
	Op   string                 `json:"op"`
	Args map[string]interface{} `json:"args"`
}

type Response struct {
	Types       Types
	ProductCode types.ProductCode
	Board       markets.Orderbook
	Ticker      markets.Ticker
	Executions  []markets.Trade

	ChildOrderEvent  []fills.ChildOrderFill
	ParentOrderEvent []fills.ParentOrderFill

	Results error
}

func subscribe(conn *websocket.Conn, channels, symbols []string) (err error) {
	var requests []request
	if symbols != nil {
		for i := range channels {
			for j := range symbols {
				requests = append(requests, request{
					Jsonrpc: "2.0",
					Method:  "subscribe",
					Params: map[string]interface{}{
						"channel": fmt.Sprintf("%s_%s", channels[i], symbols[j]),
					},
					ID: 1,
				})
			}
		}
	} else {
		for i := range channels {
			requests = append(requests, request{
				Jsonrpc: "2.0",
				Method:  "subscribe",
				Params: map[string]interface{}{
					"channel": channels[i],
				},
				ID: 1,
			})
		}
	}

	fmt.Printf("%+v\n", requests)

	for i := range requests {
		if err := conn.WriteJSON(requests[i]); err != nil {
			return err
		}
	}

	return nil
}

func unsubscribe(conn *websocket.Conn, channels, symbols []string) {
	if symbols != nil {
		for i := range channels {
			for j := range symbols {
				if err := conn.WriteJSON(&request{
					Jsonrpc: "2.0",
					Method:  "unsubscribe",
					Params: map[string]interface{}{
						"channel": fmt.Sprintf("%s_%s", channels[i], symbols[j]),
					},
					ID: 1,
				}); err != nil {
					fmt.Printf("%+v\n", err)
				}
			}
		}
	} else {
		for i := range channels {
			if err := conn.WriteJSON(&request{
				Jsonrpc: "2.0",
				Method:  "unsubscribe",
				Params: map[string]interface{}{
					"channel": channels[i],
				},
				ID: 1,
			}); err != nil {
				fmt.Printf("%+v\n", err)
			}
		}
	}
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, l *log.Logger) error {
	if l == nil {
		l = log.New(os.Stdout, "bitflyer websocket", log.Llongfile)
	}

RECONNECT:
	conn, _, err := websocket.DefaultDialer.Dial(USE1, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = subscribe(conn, channels, symbols)
	if err != nil {
		log.Fatalf("disconnect %v", err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		defer conn.Close()
		defer unsubscribe(conn, channels, symbols)

		for {
			var res Response
			_, msg, err := conn.ReadMessage()
			if err != nil {
				l.Printf("[ERROR]: msg error: %+v", err)
				res.Types = ERROR
				res.Results = fmt.Errorf("%v", err)
				ch <- res
				return fmt.Errorf("can't receive error: %v", err)
			}

			name, err := jsonparser.GetString(msg, "params", "channel")
			if err != nil {
				l.Printf("[ERROR]: channel error: %+v", err)
				res.Types = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			data, _, _, err := jsonparser.Get(msg, "params", "message")
			if err != nil {
				l.Printf("[ERROR]: message error: %+v", err)
				res.Types = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			var w Response

			switch {
			case strings.HasPrefix(name, "lightning_board_snapshot_"):
				w.ProductCode = types.ProductCode(name[len("lightning_board_snapshot_"):])
				w.Types = BOARD
				if err := json.Unmarshal(data, &w.Board); err != nil {
					l.Printf("[WARN]: cant unmarshal board %+v", err)
				}

			case strings.HasPrefix(name, "lightning_board_"):
				w.ProductCode = types.ProductCode(name[len("lightning_board_"):])
				w.Types = DIFF_BOARDS
				if err := json.Unmarshal(data, &w.Board); err != nil {
					l.Printf("[WARN]: cant unmarshal diff board %+v", err)
				}

			case strings.HasPrefix(name, "lightning_ticker_"):
				w.ProductCode = types.ProductCode(name[len("lightning_ticker_"):])
				w.Types = TICKER
				if err := json.Unmarshal(data, &w.Ticker); err != nil {
					l.Printf("[WARN]: cant unmarshal ticker %+v", err)
				}

			case strings.HasPrefix(name, "lightning_executions_"):
				w.ProductCode = types.ProductCode(name[len("lightning_executions_"):])
				w.Types = EXECUTIONS
				if err := json.Unmarshal(data, &w.Executions); err != nil {
					l.Printf("[WARN]: cant unmarshal executions %+v", err)
				}

			default:
				w.Types = UNDEFINED
				w.Results = fmt.Errorf("%v", string(msg))
			}

			select { // 外部からの停止
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// log.Debugf("recieve to send time: %v\n", time.Now().Sub(start))
			ch <- w
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("%v", err)
	}

	// 明示的 Unsubscribed
	// context.cancel()された場合は
	unsubscribe(conn, channels, symbols)

	// Maintenanceならば待機
	// Maintenanceでなければ、即再接続
	if isMentenance() {
		for {
			if !isMentenance() {
				break
			}
			time.Sleep(time.Second)
		}
	}

	goto RECONNECT
}

func requestsForPrivate(conn *websocket.Conn, key, secret string) error {
	now, nonce, sign := WsParamForPrivate(secret)
	req := &request{
		Jsonrpc: "2.0",
		Method:  "auth",
		Params: map[string]interface{}{
			"api_key":   key,
			"timestamp": now,
			"nonce":     nonce,
			"signature": sign,
		},
		ID: now,
	}

	if err := conn.WriteJSON(req); err != nil {
		return err
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	isSuccess, _ := jsonparser.GetBoolean(msg, "result")
	if !isSuccess { // read channel return, if result  false
		return err
	}
	fmt.Printf("private channel connect success: %t\n", isSuccess)

	return nil
}

func ConnectForPrivate(ctx context.Context, ch chan Response, key, secret string, channels []string, l *log.Logger) {
	if l == nil {
		l = log.New(os.Stdout, "ftx websocket", log.Llongfile)
	}

RECONNECT:
	conn, _, err := websocket.DefaultDialer.Dial(USE1, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := requestsForPrivate(conn, key, secret); err != nil {
		log.Fatalf("cant connect to private %v", err)
	}

	err = subscribe(conn, channels, nil)
	if err != nil {
		log.Fatalf("disconnect %v", err)
	}
	defer unsubscribe(conn, channels, nil)

	var eg errgroup.Group
	eg.Go(func() error {
		defer conn.Close()
		defer unsubscribe(conn, channels, nil)

		for {
			var res Response
			_, msg, err := conn.ReadMessage()
			if err != nil {
				l.Printf("[ERROR]: msg error: %+v", err)
				res.Types = ERROR
				res.Results = fmt.Errorf("%v", err)
				ch <- res
				return fmt.Errorf("can't receive error: %v", err)
			}

			name, err := jsonparser.GetString(msg, "params", "channel")
			if err != nil {
				l.Printf("[ERROR]: channel error: %+v", string(msg))
				res.Types = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			data, _, _, err := jsonparser.Get(msg, "params", "message")
			if err != nil {
				l.Printf("[ERROR]: message error: %+v", string(msg))
				res.Types = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
				continue
			}

			var w Response

			switch {
			case strings.HasPrefix(name, "child_order_events"):
				w.Types = CHILD_ORDERS
				if err := json.Unmarshal(data, &w.ChildOrderEvent); err != nil {
					l.Printf("[WARN]: cant unmarshal child_order_events %+v", err)
					continue
				}

			case strings.HasPrefix(name, "parent_order_events"):
				w.Types = PARENT_ORDERS
				if err := json.Unmarshal(data, &w.ParentOrderEvent); err != nil {
					l.Printf("[WARN]: cant unmarshal parent_order_events %+v", err)
					continue
				}

			default:
				w.Types = UNDEFINED
				w.Results = fmt.Errorf("%v", string(msg))
			}

			select { // 外部からの停止
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// log.Debugf("recieve to send time: %v\n", time.Now().Sub(start))
			ch <- w
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("%v", err)
	}

	goto RECONNECT
}

func isMentenance() bool {
	// ServerTimeを考慮し、UTC基準に
	hour := time.Now().UTC().Hour()
	if hour != 19 {
		return false
	}

	if 12 < time.Now().Minute() { // メンテナンス以外
		return false
	}
	return true
}

func WsParamForPrivate(sercret string) (now int, nonce, sign string) {
	mac := hmac.New(sha256.New, []byte(sercret))

	t := time.Now().UTC()
	rand.Seed(t.UnixNano())

	now = int(t.Unix())
	nonce = fmt.Sprintf("%d", rand.Int())

	mac.Write([]byte(fmt.Sprintf("%d%s", now, nonce)))

	sign = hex.EncodeToString(mac.Sum(nil))
	return now, nonce, sign
}
