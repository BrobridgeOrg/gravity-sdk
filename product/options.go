package product

type Options struct {
	Endpoint string
	Domain   string
	Verbose  bool
}

func NewOptions() *Options {
	options := &Options{
		Endpoint: "default",
		Domain:   "gravity",
		Verbose:  false,
	}

	return options
}
