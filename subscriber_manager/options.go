package subscriber_manager

type Options struct {
	Endpoint string
	Domain   string
	AppID    string
	AppKey   string
	Verbose  bool
}

func NewOptions() *Options {
	return &Options{
		Endpoint: "default",
		Domain:   "gravity",
		AppID:    "",
		AppKey:   "",
		Verbose:  false,
	}
}
