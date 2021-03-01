package realtime

import "log"

type Configuration struct {
	l   *log.Logger
	url string
}

type Option func(*Configuration)

func Config(ops ...Option) *Configuration {
	cfg := Configuration{
		l:   nil,
		url: "wss://ws.lightstream.bitflyer.com/json-rpc",
	}
	for _, option := range ops {
		option(&cfg)
	}
	return &cfg
}
