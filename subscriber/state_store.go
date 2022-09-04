package subscriber

type StateStore interface {
	GetPipelineState(uint64) (PipelineState, error)
	GetPipelines() []uint64
}

type PipelineState interface {
	GetLastSequence() uint64
	UpdateLastSequence(uint64) error
	Flush() error
}

// Dummy for testing
func NewStateStoreDummy() StateStore {
	return &StateStoreDummy{}
}

type StateStoreDummy struct {
}

func (ssd *StateStoreDummy) GetPipelineState(id uint64) (PipelineState, error) {
	return &PipelineStateDummy{}, nil
}

func (ssd *StateStoreDummy) GetPipelines() []uint64 {
	return []uint64{0}
}

type PipelineStateDummy struct {
	lastSeq uint64
}

func (psd *PipelineStateDummy) GetLastSequence() uint64 {
	return psd.lastSeq
}

func (psd *PipelineStateDummy) UpdateLastSequence(lastSeq uint64) error {
	psd.lastSeq = lastSeq
	return nil
}

func (psd *PipelineStateDummy) Flush() error {
	return nil
}
