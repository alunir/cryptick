package realtime

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/goccy/go-json"

	"github.com/alunir/cryptick/eodhistoricaldata/types/markets"
	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

const (
	UNDEFINED = iota
	ERROR
	US_TRADE
	US_QUOTE
	US_INDICES
	FOREX
	CRYPTO
)

type request struct {
	Action  string `json:"action"`
	Symbols string `json:"symbols"`
}

type Response struct {
	Type    int
	Symbol  string
	UsTrade markets.UsTrade
	UsQuote markets.UsQuote
	UsIndex markets.UsIndex
	Forex   markets.Forex
	Crypto  markets.Crypto
	Results error
}

func subscribe(ctx context.Context, conn *websocket.Conn, symbols []string) error {
	message := strings.Join(symbols, ",")
	if val, err := json.Marshal(request{
		Action:  "subscribe",
		Symbols: message,
	}); err != nil {
		return err
	} else if err := conn.Write(ctx, websocket.MessageBinary, val); err != nil {
		return err
	}
	return nil
}

func unsubscribe(ctx context.Context, conn *websocket.Conn, symbols []string) error {
	message := strings.Join(symbols, ",")
	if val, err := json.Marshal(request{
		Action:  "unsubscribe",
		Symbols: message,
	}); err != nil {
		return err
	} else if err := conn.Write(ctx, websocket.MessageBinary, val); err != nil {
		return err
	}
	return nil
}

func Connect(ctx context.Context, ch chan Response, symbols []string, cfg *Configuration) error {
	if cfg.l == nil {
		cfg.l = log.New(os.Stdout, "eodhistoricaldata websocket", log.Llongfile)
	}

	var isSubscribe bool

RECONNECT:
	isSubscribe = false
	conn, _, err := websocket.Dial(ctx, cfg.url, nil)
	if err != nil {
		cfg.l.Fatal(err)
	}

	if err := subscribe(ctx, conn, symbols); err != nil {
		cfg.l.Fatal(err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		defer conn.Close(websocket.StatusNormalClosure, "normal closure")
		defer unsubscribe(ctx, conn, symbols)

		for {
			var res Response
			messageType, msg, err := conn.Read(ctx)
			if err != nil {
				cfg.l.Printf("[ERROR]: msg error: %+v", string(msg))
				res.Type = ERROR
				res.Results = fmt.Errorf("%v", err)
				ch <- res
				return fmt.Errorf("can't receive error: %v", err)
			}

			if messageType == websocket.MessageText {
				if !isSubscribe {
					statusCode, _ := jsonparser.GetInt(msg, "status_code")
					if statusCode == 200 {
						cfg.l.Printf("[SUCCESS]: connection established")
						isSubscribe = true
						continue
					} else {
						cfg.l.Printf("[ERROR]: msg error: %+v", string(msg))
						res.Type = ERROR
						res.Results = fmt.Errorf("%v", string(msg))
						ch <- res
						return fmt.Errorf("can't receive error: %v", string(msg))
					}
				} else {
					switch cfg.group {
					case GROUP_US_TRADE:
						res.Type = US_TRADE
						if err := json.Unmarshal(msg, &res.UsTrade); err != nil {
							cfg.l.Printf("[WARN]: cant unmarshal us_trade %+v", err)
							continue
						}
					case GROUP_US_QUOTE:
						res.Type = US_QUOTE
						if err := json.Unmarshal(msg, &res.UsQuote); err != nil {
							cfg.l.Printf("[WARN]: cant unmarshal us_quote %+v", err)
							continue
						}
					case GROUP_US_INDICES:
						res.Type = US_INDICES
						if err := json.Unmarshal(msg, &res.UsIndex); err != nil {
							cfg.l.Printf("[WARN]: cant unmarshal us_index %+v", err)
							continue
						}
					case GROUP_FOREX:
						res.Type = FOREX
						if err := json.Unmarshal(msg, &res.Forex); err != nil {
							cfg.l.Printf("[WARN]: cant unmarshal forex %+v", err)
							continue
						}
					case GROUP_CRYPTO:
						res.Type = CRYPTO
						if err := json.Unmarshal(msg, &res.Crypto); err != nil {
							cfg.l.Printf("[WARN]: cant unmarshal crypto %+v", err)
							continue
						}
					default:
						res.Type = UNDEFINED
						res.Results = fmt.Errorf("%v", string(msg))
					}
				}
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
