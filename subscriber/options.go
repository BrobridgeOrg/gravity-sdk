package subscriber

type Options struct {
	BufferSize int
	ChunkSize  int
	Verbose    bool
}

func NewOptions() *Options {
	return &Options{
		BufferSize: 20480,
		ChunkSize:  2048,
		Verbose:    false,
	}
}
