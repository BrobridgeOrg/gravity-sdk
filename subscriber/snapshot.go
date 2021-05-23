package subscriber

import uuid "github.com/satori/go.uuid"

type Snapshot struct {
	pipeline *Pipeline
	id       string
	count    uint64
	total    uint64
	lastKey  string
}

func NewSnapshot(pipeline *Pipeline) *Snapshot {
	return &Snapshot{
		pipeline: pipeline,
		id:       uuid.NewV1().String(),
	}
}

func (snapshot *Snapshot) Create() error {
	return nil
}

func (snapshot *Snapshot) Pull() error {
	return nil
}
