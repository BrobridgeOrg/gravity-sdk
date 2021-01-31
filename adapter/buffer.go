package adapter

type RequestBuffer struct {
	size   int
	queue  []*Request
	output chan []*Request
}

func NewRequestBuffer(size int) *RequestBuffer {
	return &RequestBuffer{
		size:   size,
		queue:  make([]*Request, 0, size),
		output: make(chan []*Request, 1024),
	}
}

func (rb *RequestBuffer) SetSize(size int) {
	rb.size = size
	rb.queue = make([]*Request, 0, size)
}

func (rb *RequestBuffer) Push(req *Request) {

	rb.queue = append(rb.queue, req)

	if len(rb.queue) == rb.size {
		rb.Flush()
	}
}

func (rb *RequestBuffer) Flush() error {

	if len(rb.queue) == 0 {
		return nil
	}

	queue := rb.queue
	rb.queue = make([]*Request, 0, rb.size)
	rb.output <- queue

	return nil
}
