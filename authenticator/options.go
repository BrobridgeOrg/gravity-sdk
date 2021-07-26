package authenticator

import "fmt"

type Options struct {
	Endpoint  string
	Domain    string
	Verbose   bool
	Channel   string
	AppID     string
	AccessKey string
}

func NewOptions() *Options {
	return &Options{
		Endpoint:  "default",
		Domain:    "gravity",
		Verbose:   false,
		Channel:   "",
		AppID:     "gravity",
		AccessKey: "",
	}
}

func (options *Options) GetChannel() string {

	if len(options.Channel) == 0 {
		return fmt.Sprintf("%s.auth", options.Domain)
	}

	return options.Channel
}
