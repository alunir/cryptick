package realtime

import "log"

type Configuration struct {
	l      *log.Logger
	isTest bool
	url    string
}

type Option func(*Configuration)

func TestNet() Option {
	return func(c *Configuration) {
		c.isTest = true
		c.url = "wss://testapp.deribit.com/ws/api/v2/"
	}
}

func Config(ops ...Option) *Configuration {
	cfg := Configuration{
		l:      nil,
		isTest: false,
		url:    "wss://www.deribit.com/ws/api/v2/",
	}
	for _, option := range ops {
		option(&cfg)
	}
	return &cfg
}
