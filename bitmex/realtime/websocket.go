package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alunir/cryptick/bitmex/types/markets"
	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const (
	// Bitmex websocket op
	BitmexWSInstrument     = "instrument"
	BitmexWSOrderBookL2_25 = "orderBookL2_25"
	BitmexWSOrderBookL2    = "orderBookL2"
	BitmexWSOrderBook10    = "orderBook10"
	BitmexWSQuote          = "quote"
	BitmexWSTrade          = "trade"
	BitmexWSTradeBin1m     = "tradeBin1m"
	BitmexWSTradeBin5m     = "tradeBin5m"
	BitmexWSTradeBin1h     = "tradeBin1h"
	BitmexWSTradeBin1d     = "tradeBin1d"
	BitmexWSExecution      = "execution"
	BitmexWSOrder          = "order"
	BitmexWSMargin         = "margin"
	BitmexWSPosition       = "position"
	BitmexWSWallet         = "wallet"
)

const (
	UNDEFINED = iota
	ERROR
	TICKER
	TRADES
	ORDERBOOK
	ORDERBOOK_L2
	ORDERS
	FILLS
)

type request struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type Response struct {
	Type        int
	Ticker      []markets.Ticker
	Trades      []markets.Trade
	OrderbookL2 []markets.OrderBookL2
	Orderbook   []markets.OrderBook10
	// TODO: Not implemented yet
	// Orders    orders.Order
	// Fills     fills.Fill
	Results error
}

func subscribe(conn *websocket.Conn, channels, symbols []string) error {
	var message []interface{}
	for i := range channels {
		for j := range symbols {
			message = append(message, fmt.Sprintf("%v:%v", channels[i], symbols[j]))
		}
	}
	if err := conn.WriteJSON(&request{
		Op:   "subscribe",
		Args: message,
	}); err != nil {
		return err
	}
	return nil
}

func unsubscribe(conn *websocket.Conn, channels, symbols []string) error {
	var message []interface{}
	for i := range channels {
		for j := range symbols {
			message = append(message, fmt.Sprintf("%v:%v", channels[i], symbols[j]))
		}
	}
	if err := conn.WriteJSON(&request{
		Op:   "unsubscribe",
		Args: message,
	}); err != nil {
		return err
	}
	return nil
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, cfg *Configuration) error {
	if cfg.l == nil {
		cfg.l = log.New(os.Stdout, "bitmex websocket", log.Llongfile)
	}

RECONNECT:
	conn, _, err := websocket.DefaultDialer.Dial(cfg.url, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := subscribe(conn, channels, symbols); err != nil {
		log.Fatal(err)
	}

	go ping(conn)

	var eg errgroup.Group
	eg.Go(func() error {
		defer conn.Close()
		defer unsubscribe(conn, channels, symbols)

		for {
			var res Response
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				cfg.l.Printf("[ERROR]: msg error: %+v", err)
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", err)
				ch <- res
				return fmt.Errorf("can't receive error: %v", err)
			}

			if messageType == websocket.TextMessage {
				if string(msg) == "pong" {
					continue
				}
			}

			success, err := jsonparser.GetBoolean(msg, "success")
			if success {
				cfg.l.Printf("[SUCCESS]: %+v", string(msg))
				continue
			}

			channel, err := jsonparser.GetString(msg, "table")
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

			switch channel {
			case BitmexWSQuote:
				res.Type = TICKER
				if err := json.Unmarshal(data, &res.Ticker); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal ticker %+v", err)
					continue
				}
			case BitmexWSTrade:
				res.Type = TRADES
				if err := json.Unmarshal(data, &res.Trades); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal trades %+v", err)
					continue
				}
			case BitmexWSOrderBookL2:
				res.Type = ORDERBOOK_L2
				if err := json.Unmarshal(data, &res.OrderbookL2); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal orderbookL2 %+v", err)
					continue
				}
			case BitmexWSOrderBookL2_25:
				res.Type = ORDERBOOK_L2
				if err := json.Unmarshal(data, &res.OrderbookL2); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal orderbookL2_25 %+v", err)
					continue
				}
			case BitmexWSOrderBook10:
				res.Type = ORDERBOOK
				if err := json.Unmarshal(data, &res.Orderbook); err != nil {
					cfg.l.Printf("[WARN]: cant unmarshal orderbook10 %+v", err)
					continue
				}
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

func ping(conn *websocket.Conn) (err error) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				goto EXIT
			}
		}
	}
EXIT:
	return err
}
