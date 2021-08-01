package eventstore

import "github.com/BrobridgeOrg/gravity-sdk/core/keyring"

type Options struct {
	Endpoint string
	Domain   string
	Verbose  bool
	Key      *keyring.KeyInfo
}

func NewOptions() *Options {
	options := &Options{
		Endpoint: "default",
		Domain:   "gravity",
		Verbose:  false,
		Key:      keyring.NewKey("gravity", ""),
	}

	return options
}
