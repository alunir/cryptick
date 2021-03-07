package realtime

import "log"

type category string
type group string
type base string

var (
	TESTNET  = category("testnet")
	MAINNET1 = category("mainnet1")
	MAINNET2 = category("mainnet2")

	PUBLIC  = group("public")
	PRIVATE = group("private")

	PERPETUAL = base("perpetual")
	USDT      = base("usdt")
)

type Endpoint struct {
	url      string
	category category
	group    group
	base     base
}

var (
	PERPETUAL_TESTNET     = Endpoint{"wss://stream-testnet.bybit.com/realtime", TESTNET, PUBLIC, PERPETUAL}
	PERPETUAL_MAINNET1    = Endpoint{"wss://stream.bybit.com/realtime", MAINNET1, PUBLIC, PERPETUAL}
	PERPETUAL_MAINNET2    = Endpoint{"wss://stream.bytick.com/realtime", MAINNET2, PUBLIC, PERPETUAL}
	USDT_TESTNET_PUBLIC   = Endpoint{"wss://stream-testnet.bybit.com/realtime_public", TESTNET, PUBLIC, USDT}
	USDT_TESTNET_PRIVATE  = Endpoint{"wss://stream-testnet.bybit.com/realtime_private", TESTNET, PRIVATE, USDT}
	USDT_MAINNET1_PUBLIC  = Endpoint{"wss://stream.bybit.com/realtime_public", MAINNET1, PUBLIC, USDT}
	USDT_MAINNET2_PUBLIC  = Endpoint{"wss://stream.bytick.com/realtime_public", MAINNET2, PUBLIC, USDT}
	USDT_MAINNET1_PRIVATE = Endpoint{"wss://stream.bybit.com/realtime_private", MAINNET1, PRIVATE, USDT}
	USDT_MAINNET2_PRIVATE = Endpoint{"wss://stream.bytick.com/realtime_private", MAINNET2, PRIVATE, USDT}
)

type Configuration struct {
	l      *log.Logger
	url    string
	key    string
	secret string
}

type Option func(*Configuration)

func EndpointOption(e Endpoint) Option {
	return func(c *Configuration) {
		c.url = e.url
	}
}

func Key(key string) Option {
	return func(c *Configuration) {
		c.key = key
	}
}

func SecretKey(secret string) Option {
	return func(c *Configuration) {
		c.secret = secret
	}
}

func Config(ops ...Option) *Configuration {
	cfg := Configuration{
		l:      nil,
		url:    PERPETUAL_MAINNET1.url,
		key:    "",
		secret: "",
	}
	for _, option := range ops {
		option(&cfg)
	}
	return &cfg
}
