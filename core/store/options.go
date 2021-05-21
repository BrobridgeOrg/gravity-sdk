package store

import broton "github.com/BrobridgeOrg/broton"

type Options struct {
	StoreOptions *broton.Options
}

func NewOptions() *Options {
	return &Options{
		StoreOptions: broton.NewOptions(),
	}
}
