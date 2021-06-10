package state_store

import (
	broton "github.com/BrobridgeOrg/broton"
)

type PipelineState struct {
	stateStore *StateStore
	pipelineID uint64
	lastSeq    uint64
}

func NewPipelineState(ss *StateStore, pipelineID uint64) *PipelineState {
	return &PipelineState{
		stateStore: ss,
		pipelineID: pipelineID,
	}
}

func (ps *PipelineState) Initialize() error {

	// Loading pipeline state from store
	store := ps.stateStore.store
	value, err := store.GetUint64("pipeline", broton.Uint64ToBytes(ps.pipelineID))
	if err != nil {
		return err
	}

	ps.lastSeq = value

	return nil
}

func (ps *PipelineState) UpdateLastSequence(lastSeq uint64) error {

	store := ps.stateStore.store
	err := store.PutUint64("pipeline", broton.Uint64ToBytes(ps.pipelineID), lastSeq)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PipelineState) GetLastSequence() uint64 {
	return ps.lastSeq
}

func (ps *PipelineState) Flush() error {

	cf, err := ps.stateStore.store.GetColumnFamailyHandle("pipeline")
	if err != nil {
		return err
	}

	_, err = cf.Db.AsyncFlush()
	if err != nil {
		return err
	}

	return nil
}
