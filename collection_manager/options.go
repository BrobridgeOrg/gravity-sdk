package collection_manager

import "github.com/BrobridgeOrg/gravity-sdk/core/keyring"

type Options struct {
	Endpoint string
	Domain   string
	Key      *keyring.KeyInfo
	Verbose  bool
}

func NewOptions() *Options {
	return &Options{
		Endpoint: "default",
		Domain:   "gravity",
		Key:      keyring.NewKey("gravity", ""),
		Verbose:  false,
	}
}
