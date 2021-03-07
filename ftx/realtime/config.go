package realtime

import "log"

type Configuration struct {
	l          *log.Logger
	key        string
	secret     string
	url        string
	subaccount []string
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

func Subaccount(subaccount ...string) Option {
	return func(c *Configuration) {
		c.subaccount = subaccount
	}
}

func Config(ops ...Option) *Configuration {
	cfg := Configuration{
		l:          nil,
		key:        "",
		secret:     "",
		url:        "wss://ftx.com/ws/",
		subaccount: []string{},
	}
	for _, option := range ops {
		option(&cfg)
	}
	return &cfg
}
