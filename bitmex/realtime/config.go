package realtime

import "log"

type Configuration struct {
	l      *log.Logger
	isTest bool
	url    string
	key    string
	secret string
}

type Option func(*Configuration)

func TestNet() Option {
	return func(c *Configuration) {
		c.isTest = true
		c.url = "wss://testnet.bitmex.com/realtime"
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
		isTest: false,
		url:    "wss://www.bitmex.com/realtime",
		key:    "",
		secret: "",
	}
	for _, option := range ops {
		option(&cfg)
	}
	return &cfg
}
