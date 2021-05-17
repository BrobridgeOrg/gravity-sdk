package state_store

import (
	"sync"

	broton "github.com/BrobridgeOrg/broton"
	gravity_subscriber "github.com/BrobridgeOrg/gravity-sdk/subscriber"
)

type StateStore struct {
	options     *Options
	storeEngine *broton.Broton
	store       *broton.Store
	collections sync.Map
	pipelines   map[uint64]*PipelineState
}

func NewStateStore(options *Options) *StateStore {
	return &StateStore{
		options:   options,
		pipelines: make(map[uint64]*PipelineState),
	}
}

func (ss *StateStore) Initialize() error {

	es, err := broton.NewBroton(ss.options.StoreOptions)
	if err != nil {
		return err
	}

	ss.storeEngine = es

	// Initializing store
	store, err := es.GetStore("gravity_state_store")
	if err != nil {
		return err
	}

	ss.store = store

	// Initializing columns
	err = store.RegisterColumns([]string{"pipeline"})
	if err != nil {
		return err
	}

	return nil
}

func (ss *StateStore) createPipeline(pipelineID uint64) (gravity_subscriber.PipelineState, error) {

	pipelineState := NewPipelineState(ss, pipelineID)
	err := pipelineState.Initialize()
	if err != nil {
		return nil, err
	}

	ss.pipelines[pipelineID] = pipelineState

	return gravity_subscriber.PipelineState(pipelineState), nil
}

func (ss *StateStore) GetPipelineState(pipelineID uint64) (gravity_subscriber.PipelineState, error) {

	pipeline, ok := ss.pipelines[pipelineID]
	if !ok {
		return ss.createPipeline(pipelineID)
	}

	return gravity_subscriber.PipelineState(pipeline), nil
}

func (ss *StateStore) GetPipelines() []uint64 {

	pipelines := make([]uint64, 0, len(ss.pipelines))
	for pipelineID, _ := range ss.pipelines {
		pipelines = append(pipelines, pipelineID)
	}

	return pipelines
}
