package subscriber

type Options struct {
	WorkerCount int
	BufferSize  int
	ChunkSize   int
	Verbose     bool
	StateStore  StateStore
}

func NewOptions() *Options {
	return &Options{
		WorkerCount: 4,
		BufferSize:  20480,
		ChunkSize:   2048,
		Verbose:     false,
		StateStore:  nil,
	}
}
