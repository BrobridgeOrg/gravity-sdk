package authenticator

import (
	"fmt"

	"github.com/BrobridgeOrg/gravity-sdk/core/keyring"
)

type Options struct {
	Endpoint string
	Domain   string
	Verbose  bool
	Channel  string
	Key      *keyring.KeyInfo
}

func NewOptions() *Options {
	return &Options{
		Endpoint: "default",
		Domain:   "gravity",
		Verbose:  false,
		Channel:  "",
		Key:      keyring.NewKey("gravity", ""),
	}
}

func (options *Options) GetChannel() string {

	if len(options.Channel) == 0 {
		return fmt.Sprintf("%s.auth", options.Domain)
	}

	return options.Channel
}
