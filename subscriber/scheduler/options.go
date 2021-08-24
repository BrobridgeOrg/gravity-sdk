package scheduler

type Options struct {
	WorkerCount int
}

func NewOptions() *Options {
	return &Options{
		WorkerCount: 4,
	}
}
