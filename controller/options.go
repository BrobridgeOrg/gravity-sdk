package controller

type Options struct {
	BufferSize int
}

func NewOptions() *Options {
	return &Options{
		BufferSize: 20480,
	}
}
