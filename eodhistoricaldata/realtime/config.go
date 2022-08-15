package realtime

import (
	"log"
)

const ENDPOINT_BASE = "wss://ws.eodhistoricaldata.com/ws/"

type group string

var (
	GROUP_US_TRADE   = group("us")
	GROUP_US_QUOTE   = group("us-quote")
	GROUP_US_INDICES = group("index")
	GROUP_FOREX      = group("forex")
	GROUP_CRYPTO     = group("crypto")
)

type Endpoint struct {
	group group
}

var (
	ENDPOINT_US_TRADE   = Endpoint{GROUP_US_TRADE}
	ENDPOINT_US_QUOTE   = Endpoint{GROUP_US_QUOTE}
	ENDPOINT_US_INDICES = Endpoint{GROUP_US_INDICES}
	ENDPOINT_FOREX      = Endpoint{GROUP_FOREX}
	ENDPOINT_CRYPTO     = Endpoint{GROUP_CRYPTO}
)

type Configuration struct {
	l     *log.Logger
	group group
	url   string
	key   string
}

type Option func(*Configuration)

func EndpointGroup(e Endpoint) Option {
	return func(c *Configuration) {
		c.group = e.group
	}
}

func Key(key string) Option {
	return func(c *Configuration) {
		c.key = key
	}
}

func Config(ops ...Option) *Configuration {
	cfg := Configuration{
		l:   nil,
		key: "",
	}
	for _, option := range ops {
		option(&cfg)
	}
	if cfg.key == "" {
		panic("key is empty")
	}
	cfg.url = ENDPOINT_BASE + string(cfg.group) + "?api_token=" + cfg.key
	return &cfg
}
