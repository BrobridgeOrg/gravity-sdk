package subscriber

type StateStore interface {
	GetPipelineState(uint64) (PipelineState, error)
	GetPipelines() []uint64
}

type PipelineState interface {
	GetLastSequence() uint64
	UpdateLastSequence(uint64) error
}
