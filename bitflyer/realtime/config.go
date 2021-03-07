package realtime

import "log"

type Configuration struct {
	l      *log.Logger
	url    string
	key    string
	secret string
}

type Option func(*Configuration)

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
		url:    "wss://ws.lightstream.bitflyer.com/json-rpc",
		key:    "",
		secret: "",
	}
	for _, option := range ops {
		option(&cfg)
	}
	return &cfg
}
