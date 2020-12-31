package adapter

import "time"

type Options struct {
	PingInterval        time.Duration
	MaxPingsOutstanding int
	MaxReconnects       int
	ReconnectHandler    func()
	DisconnectHandler   func()
}

func NewOptions() *Options {
	return &Options{
		PingInterval:        10 * time.Second,
		MaxPingsOutstanding: 3,
		MaxReconnects:       -1,
		ReconnectHandler:    func() {},
		DisconnectHandler:   func() {},
	}
}
