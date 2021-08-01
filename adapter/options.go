package adapter

import (
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/core/keyring"
)

type Options struct {
	Endpoint            string
	Domain              string
	Key                 *keyring.KeyInfo
	BatchSize           int
	PingInterval        time.Duration
	MaxPingsOutstanding int
	MaxReconnects       int
	ReconnectHandler    func()
	DisconnectHandler   func()
	Verbose             bool
}

func NewOptions() *Options {
	return &Options{
		Endpoint:            "default",
		Domain:              "gravity",
		Key:                 keyring.NewKey("anonymous", ""),
		BatchSize:           1000,
		PingInterval:        10 * time.Second,
		MaxPingsOutstanding: 3,
		MaxReconnects:       -1,
		ReconnectHandler:    func() {},
		DisconnectHandler:   func() {},
		Verbose:             false,
	}
}
