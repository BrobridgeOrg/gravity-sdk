package state_store

import "github.com/BrobridgeOrg/gravity-sdk/core/store"

type Options struct {
	Core *store.Options
	Name string
}

func NewOptions() *Options {
	return &Options{
		Core: store.NewOptions(),
	}
}
