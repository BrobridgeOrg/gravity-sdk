package state_store

import (
	"sync"

	broton "github.com/BrobridgeOrg/broton"
	"github.com/BrobridgeOrg/gravity-sdk/core/store"
	gravity_subscriber "github.com/BrobridgeOrg/gravity-sdk/subscriber"
)

type StateStore struct {
	options      *Options
	gravityStore *store.Store
	store        *broton.Store
	collections  sync.Map
	pipelines    map[uint64]*PipelineState
}

func NewStateStore(options *Options) *StateStore {
	return &StateStore{
		options:   options,
		pipelines: make(map[uint64]*PipelineState),
	}
}

func NewStateStoreWithStore(gravityStore *store.Store, options *Options) *StateStore {

	ss := NewStateStore(options)
	ss.gravityStore = gravityStore
	return ss
}

func (ss *StateStore) Initialize() error {

	if ss.gravityStore == nil {
		gravityStore, err := store.NewStore(ss.options.Core)
		if err != nil {
			return err
		}

		ss.gravityStore = gravityStore
	}

	// Initializing store
	storeName := "gravity_state_store"
	if len(ss.options.Name) > 0 {
		storeName = "gravity_state_store_" + ss.options.Name
	}
	s, err := ss.gravityStore.GetEngine().GetStore(storeName)
	if err != nil {
		return err
	}

	ss.store = s

	// Initializing columns
	err = s.RegisterColumns([]string{"pipeline"})
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
