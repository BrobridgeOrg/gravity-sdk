package subscriber

type Options struct {
	Domain  string
	Verbose bool
}

func NewOptions() *Options {
	return &Options{
		Domain:  "default",
		Verbose: false,
	}
}
