package subscriber

type Options struct {
	BufferSize int
	ChunkSize  int
	Verbose    bool
	StateStore StateStore
}

func NewOptions() *Options {
	return &Options{
		BufferSize: 20480,
		ChunkSize:  2048,
		Verbose:    false,
		StateStore: nil,
	}
}
