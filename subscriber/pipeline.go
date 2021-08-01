package subscriber

import (
	"errors"
	"fmt"
	"sync"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	synchronizer_pb "github.com/BrobridgeOrg/gravity-api/service/synchronizer"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

type Pipeline struct {
	subscriber  *Subscriber
	snapshot    *Snapshot
	id          uint64
	lastSeq     uint64
	isReady     bool
	isSuspended bool
	mutex       sync.Mutex
}

func NewPipeline(subscriber *Subscriber, id uint64, lastSeq uint64) *Pipeline {
	return &Pipeline{
		subscriber: subscriber,
		id:         id,
		lastSeq:    lastSeq,
		isReady:    false,
	}
}

func (pipeline *Pipeline) UpdateLastSequence(sequence uint64) {
	pipeline.lastSeq = sequence
}

func (pipeline *Pipeline) SaveLastSequence() error {
	if pipeline.subscriber.options.StateStore == nil {
		return nil
	}

	// Write to store
	pipelineState, _ := pipeline.subscriber.options.StateStore.GetPipelineState(pipeline.id)
	return pipelineState.UpdateLastSequence(pipeline.lastSeq)
}

func (pipeline *Pipeline) getStateFromServer() (*pipeline_pb.GetStateReply, error) {

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
	}).Info("Getting pipeline states from server")

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.getState", pipeline.id)
	request := pipeline_pb.GetStateRequest{}

	msg, _ := proto.Marshal(&request)

	respData, err := pipeline.subscriber.request(channel, msg, true)
	if err != nil {
		return nil, err
	}

	var reply pipeline_pb.GetStateReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, errors.New("Forbidden")
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Error(reply.Reason)
		return nil, errors.New(reply.Reason)
	}

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"lastSeq":  reply.LastSeq,
	}).Info("latest pipeline states from server")

	return &reply, nil

}

func (pipeline *Pipeline) initialize() error {

	// Initial load is disabled
	if !pipeline.subscriber.options.InitialLoad.Enabled {
		pipeline.Ready()
		return nil
	}

	// Getting pipeline states
	state, err := pipeline.getStateFromServer()
	if err != nil {
		return err
	}

	// check states
	if pipeline.lastSeq == 0 || pipeline.lastSeq+pipeline.subscriber.options.InitialLoad.OmittedCount < state.LastSeq {

		//TODO: It should be improved to compare number of records with pipeline lastSeq
		//TODO: Truncate

		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Info("Preparing snapshot")

		// Initializing snapshot
		snapshot := NewSnapshot(pipeline)
		pipeline.snapshot = snapshot

		// Using snapshot last sequence to initialize pipeline
		pipeline.lastSeq = state.LastSeq

		// Create snapshot
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Info("Creating snapshot")

		// Create a new snapshot
		err := pipeline.snapshot.Create()
		if err != nil {
			return err
		}

		return nil
	}

	// Ready
	pipeline.Ready()

	return nil
}

func (pipeline *Pipeline) performInitialLoad() error {

	// Pull data from snapshot
	count, err := pipeline.snapshot.Pull()
	if err != nil {
		return err
	}

	// Need to wait
	if count > 0 {
		return nil
	}

	if pipeline.snapshot.isCompleted {

		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Info("Initial load was done")

		// Write last sequence of snapshot to store
		if err := pipeline.SaveLastSequence(); err != nil {
			log.WithFields(logrus.Fields{
				"pipeline": pipeline.id,
			}).Error(err)
		}

		pipeline.snapshot.Close()

		pipeline.Ready()
		pipeline.Idle()

		return nil
	}

	pipeline.Idle()

	return nil
}

func (pipeline *Pipeline) fetch() error {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.fetch", pipeline.id)

	request := synchronizer_pb.PipelineFetchRequest{
		SubscriberID: pipeline.subscriber.id,
		PipelineID:   pipeline.id,
		StartAt:      pipeline.lastSeq,
		Offset:       1,
		Count:        int64(pipeline.subscriber.options.ChunkSize),
	}

	if pipeline.lastSeq == 0 {
		request.Offset = 0
	}

	msg, _ := proto.Marshal(&request)

	respData, err := pipeline.subscriber.request(channel, msg, true)
	if err != nil {
		return fmt.Errorf("Failed to fetch: %v", err)
	}

	var reply synchronizer_pb.PipelineFetchReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return fmt.Errorf("Failed to fetch: %v", err)
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to fetch: %s", reply.Reason)
		return errors.New(reply.Reason)
	}

	// No more event so pipeline should be suspended
	if reply.Count == 0 {

		// Trying to suspend to prevent infinite loop
		pipeline.isSuspended = true
		return nil
	}

	pipeline.UpdateLastSequence(reply.LastSeq)

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"lastSeq":  reply.LastSeq,
		"count":    reply.Count,
	}).Info("-> Fetching event chunk")

	return nil
}

func (pipeline *Pipeline) Pull() error {

	pipeline.mutex.Lock()
	defer pipeline.mutex.Unlock()

	if !pipeline.isReady {

		// Initial load is disabled
		if !pipeline.subscriber.options.InitialLoad.Enabled {
			return nil
		}

		// fetching data for initial load
		err := pipeline.performInitialLoad()
		if err != nil {
			return err
		}

		return nil
	}

	return pipeline.fetch()
}

func (pipeline *Pipeline) Suspend() bool {

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"lastSeq":  pipeline.lastSeq,
	}).Info("<- Suspending pipeline")

	pipeline.isSuspended = false

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.suspend", pipeline.id)
	request := pipeline_pb.SuspendRequest{
		SubscriberID: pipeline.subscriber.id,
		Sequence:     pipeline.lastSeq,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := pipeline.subscriber.request(channel, msg, true)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to suspend: %v", err)
		return false
	}

	var reply pipeline_pb.SuspendReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to suspend: %v", err)
		return false
	}

	if len(reply.Reason) > 0 {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to suspend: %s", reply.Reason)
		return false
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Warn("Disallow suscriber to suspend")
		return false
	}

	pipeline.isSuspended = true

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
	}).Warn("pipeline suspended")

	// Force to flush
	pipelineState, _ := pipeline.subscriber.options.StateStore.GetPipelineState(pipeline.id)
	err = pipelineState.Flush()
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
			"lastSeq":  pipeline.lastSeq,
		}).Errorf("Failed to flush state store: %s", err)
	}

	return true
}

func (pipeline *Pipeline) Idle() {
	pipeline.subscriber.ReleasePipeline(pipeline.id)
}

func (pipeline *Pipeline) Ready() {

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"lastSeq":  pipeline.lastSeq,
	}).Info("Pipeline is ready")

	pipeline.isReady = true
}

func (pipeline *Pipeline) Awake() {
	pipeline.isSuspended = false
}
