package adapter

import "time"

type Options struct {
	BatchSize           int
	PingInterval        time.Duration
	MaxPingsOutstanding int
	MaxReconnects       int
	ReconnectHandler    func()
	DisconnectHandler   func()
}

func NewOptions() *Options {
	return &Options{
		BatchSize:           1000,
		PingInterval:        10 * time.Second,
		MaxPingsOutstanding: 3,
		MaxReconnects:       -1,
		ReconnectHandler:    func() {},
		DisconnectHandler:   func() {},
	}
}
