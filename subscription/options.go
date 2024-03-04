package subscription

type Options struct {
	Domain  string
	Verbose bool
}

func NewOptions() *Options {
	options := &Options{
		Domain:  "gravity",
		Verbose: false,
	}

	return options
}
