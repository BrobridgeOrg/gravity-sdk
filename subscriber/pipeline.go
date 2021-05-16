package subscriber

type Pipeline struct {
	id      uint64
	lastSeq uint64
}

func NewPipeline(id uint64, lastSeq uint64) *Pipeline {
	return &Pipeline{
		id:      id,
		lastSeq: lastSeq,
	}
}
