package core

import "time"

type Options struct {
	Token                  string
	PingInterval           time.Duration
	MaxPingsOutstanding    int
	MaxReconnects          int
	ReconnectHandler       func()
	DisconnectHandler      func()
	PublishAsyncMaxPending int
}

func NewOptions() *Options {
	return &Options{
		PingInterval:           10 * time.Second,
		MaxPingsOutstanding:    3,
		MaxReconnects:          -1,
		ReconnectHandler:       func() {},
		DisconnectHandler:      func() {},
		PublishAsyncMaxPending: 10240,
	}
}
