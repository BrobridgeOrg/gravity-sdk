package adapter

import (
	"time"
)

type Compression int32

const (
	NoCompression Compression = iota
	S2Compression
)

type Options struct {
	Domain              string
	BatchSize           int
	PingInterval        time.Duration
	MaxPingsOutstanding int
	MaxReconnects       int
	ReconnectHandler    func()
	DisconnectHandler   func()
	Compression         Compression
	Verbose             bool
}

func NewOptions() *Options {
	return &Options{
		Domain:              "defualt",
		BatchSize:           1000,
		PingInterval:        10 * time.Second,
		MaxPingsOutstanding: 3,
		MaxReconnects:       -1,
		ReconnectHandler:    func() {},
		DisconnectHandler:   func() {},
		Compression:         NoCompression,
		Verbose:             false,
	}
}
