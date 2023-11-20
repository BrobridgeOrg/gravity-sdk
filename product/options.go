package product

type Options struct {
	Domain  string // Domain specifies the target domain or environment the options apply to.
	Verbose bool   // Verbose indicates whether verbose or detailed logging should be enabled.
}

func NewOptions() *Options {
	options := &Options{
		Domain:  "gravity",
		Verbose: false,
	}

	return options
}
