package realtime

import "log"

type Configuration struct {
	l          *log.Logger
	url        string
	subaccount []string
}

type Option func(*Configuration)

func Subaccount(subaccount ...string) Option {
	return func(c *Configuration) {
		c.subaccount = subaccount
	}
}

func Config(ops ...Option) *Configuration {
	cfg := Configuration{
		l:          nil,
		url:        "wss://ftx.com/ws/",
		subaccount: []string{},
	}
	for _, option := range ops {
		option(&cfg)
	}
	return &cfg
}
