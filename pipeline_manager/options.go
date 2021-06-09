package pipeline_manager

type Options struct {
	Endpoint string
	Domain   string
	Verbose  bool
}

func NewOptions() *Options {
	return &Options{
		Endpoint: "default",
		Domain:   "gravity",
		Verbose:  false,
	}
}
