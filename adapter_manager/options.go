package adapter_manager

type Options struct {
	Verbose bool
}

func NewOptions() *Options {
	return &Options{
		Verbose: false,
	}
}
