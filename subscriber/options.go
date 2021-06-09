package subscriber

type Options struct {
	Endpoint    string
	Domain      string
	WorkerCount int
	BufferSize  int
	ChunkSize   int
	Verbose     bool
	StateStore  StateStore
	InitialLoad InitialLoadOptions
}

type InitialLoadOptions struct {
	Enabled      bool
	OmittedCount uint64
}

func NewOptions() *Options {
	options := &Options{
		Endpoint:    "default",
		Domain:      "gravity",
		WorkerCount: 4,
		BufferSize:  20480,
		ChunkSize:   2048,
		Verbose:     false,
		StateStore:  nil,
	}

	options.InitialLoad.Enabled = false
	options.InitialLoad.OmittedCount = 100000

	return options
}
