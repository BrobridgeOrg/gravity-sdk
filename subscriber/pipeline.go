package subscriber

import (
	"errors"
	"fmt"
	"sync"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

type PipelineStatus int32

const (
	PIPELINE_STATUS_OPEN PipelineStatus = iota
	PIPELINE_STATUS_HALF_OPEN
	PIPELINE_STATUS_CLOSE
)

var PipelineStatusNames = map[PipelineStatus]string{
	0: "Open",
	1: "Half Open",
	2: "Close",
}

type Pipeline struct {
	subscriber  *Subscriber
	snapshot    *Snapshot
	id          uint64
	lastSeq     uint64
	updatedSeq  uint64
	isReady     bool
	isSuspended bool
	status      PipelineStatus
	failure     int
	mutex       sync.Mutex
}

func NewPipeline(subscriber *Subscriber, id uint64, lastSeq uint64) *Pipeline {
	return &Pipeline{
		subscriber: subscriber,
		id:         id,
		lastSeq:    lastSeq,
		updatedSeq: lastSeq,
		isReady:    false,
		status:     PIPELINE_STATUS_OPEN,
	}
}

func (pipeline *Pipeline) Initialize() error {

	// Getting pipeline states
	state, err := pipeline.getStateFromServer()
	if err != nil {
		return err
	}

	// Initial load is disabled
	if !pipeline.subscriber.options.InitialLoad.Enabled {

		if pipeline.lastSeq == 0 {
			// Using last sequence to initialize pipeline
			pipeline.SetUpdatedSequence(state.LastSeq)
		}

		pipeline.lastSeq = state.LastSeq
		pipeline.isReady = true
		return nil
	}

	if pipeline.subscriber.options.InitialLoad.Mode == "snapshot" {

		// check states for snapshot
		if pipeline.lastSeq == 0 || pipeline.lastSeq+pipeline.subscriber.options.InitialLoad.OmittedCount < state.LastSeq {

			//TODO: It should be improved to compare number of records with pipeline lastS else 		//TODO: Truncate

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
	}

	// Ready
	pipeline.isReady = true

	return nil
}

func (pipeline *Pipeline) UpdateLastSequence(sequence uint64) {
	pipeline.lastSeq = sequence
}

func (pipeline *Pipeline) SetUpdatedSequence(updatedSeq uint64) error {

	if pipeline.updatedSeq >= updatedSeq {
		return nil
	}

	pipeline.updatedSeq = updatedSeq

	if pipeline.subscriber.options.StateStore == nil {
		return nil
	}

	// Write to store
	pipelineState, _ := pipeline.subscriber.options.StateStore.GetPipelineState(pipeline.id)

	return pipelineState.UpdateLastSequence(updatedSeq)
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

func (pipeline *Pipeline) pull() ([][]byte, error) {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.pullEvents", pipeline.id)

	request := pipeline_pb.PullEventsRequest{
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
		return nil, fmt.Errorf("Failed to fetch: %v", err)
	}

	var reply pipeline_pb.PullEventsReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch: %v", err)
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to fetch: %s", reply.Reason)
		return nil, errors.New(reply.Reason)
	}

	// No more event so pipeline should be suspended
	if len(reply.Events) == 0 {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
			"curSeq":   pipeline.lastSeq,
		}).Info("No more event")
		return nil, nil
	}

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"origSeq":  pipeline.lastSeq,
		"count":    len(reply.Events),
	}).Info("-> Fetching event chunk")

	// Update last sequence
	pipeline.lastSeq = reply.LastSeq

	return reply.Events, nil
}

func (pipeline *Pipeline) suspend() error {

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"lastSeq":  pipeline.lastSeq,
	}).Info("<- Suspending pipeline")

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
		return err
	}

	var reply pipeline_pb.SuspendReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to suspend: %v", err)
		return err
	}

	if len(reply.Reason) > 0 {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to suspend: %s", reply.Reason)
		return errors.New(reply.Reason)
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Warn("Disallow subscriber to suspend")
		return errors.New("Disallow subscriber to suspend")
	}

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

	return nil
}

func (pipeline *Pipeline) awake() error {

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"lastSeq":  pipeline.lastSeq,
	}).Info("<- Suspending pipeline")

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.awake", pipeline.id)
	request := pipeline_pb.AwakeRequest{
		SubscriberID: pipeline.subscriber.id,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := pipeline.subscriber.request(channel, msg, true)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to awake: %v", err)
		return err
	}

	var reply pipeline_pb.AwakeReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to awake: %v", err)
		return err
	}

	if len(reply.Reason) > 0 {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Errorf("Failed to awake: %s", reply.Reason)
		return errors.New(reply.Reason)
	}

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
	}).Warn("pipeline is back")

	return nil
}

func (pipeline *Pipeline) dispatchMessages(msgType MessageType, privData interface{}, records [][]byte) error {

	for _, record := range records {

		msg := messagePool.Get().(*Message)
		msg.Subscription = pipeline.subscriber.subscription
		msg.Pipeline = pipeline
		msg.Type = msgType

		switch msg.Type {
		case MESSAGE_TYPE_SNAPSHOT:
			msg.Payload = &SnapshotEvent{
				PipelineID: pipeline.id,
				Collection: privData.(string),
				RawData:    record,
			}
			msg.Callback = nil

		case MESSAGE_TYPE_EVENT:
			msg.Payload = &DataEvent{
				PipelineID: pipeline.id,
				RawData:    record,
			}

			msg.Callback = func(msg *Message) {

				// Update sequence number to state store
				event := msg.Payload.(*DataEvent)
				err := pipeline.SetUpdatedSequence(event.Sequence)
				if err != nil {
					log.Errorf("Failed to write sequence to state store: %v", err)
				}
			}
		}

		pipeline.subscriber.subscription.Push(msg)
	}

	return nil
}

func (pipeline *Pipeline) Awake() error {

	pipeline.mutex.Lock()
	defer pipeline.mutex.Unlock()

	if !pipeline.isReady {

		// Pull data from snapshot
		collection, records, err := pipeline.snapshot.Pull()
		if err != nil {
			pipeline.Fail()
			return err
		}

		if records != nil {
			// Dispatch
			pipeline.dispatchMessages(MESSAGE_TYPE_SNAPSHOT, collection, records)
			pipeline.Succeed()
		}

		if pipeline.snapshot.isCompleted {
			pipeline.isReady = true
			pipeline.SetUpdatedSequence(pipeline.lastSeq)
			pipeline.snapshot.Close()
			log.WithFields(logrus.Fields{
				"pipeline": pipeline.id,
				"lastSeq":  pipeline.lastSeq,
			}).Warn("Initialized pipeline")
			return nil
		}

		return nil
	}

	// Pull events
	events, err := pipeline.pull()
	if err != nil {
		pipeline.Fail()
		return err
	}

	if events == nil {
		// No events
		pipeline.Fail()
		return nil
	}

	// Dispatch
	pipeline.dispatchMessages(MESSAGE_TYPE_EVENT, nil, events)
	pipeline.Succeed()

	return nil
}

func (pipeline *Pipeline) Succeed() {

	if pipeline.failure == 0 {
		return
	}

	pipeline.failure--
	if pipeline.failure == 0 {
		pipeline.status = PIPELINE_STATUS_OPEN
	} else if pipeline.status == PIPELINE_STATUS_CLOSE {
		// Trying to re-open
		pipeline.status = PIPELINE_STATUS_HALF_OPEN
	}

	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
	}).Infof("Switch status to %s", PipelineStatusNames[pipeline.status])
}

func (pipeline *Pipeline) Fail() {

	if pipeline.status == PIPELINE_STATUS_HALF_OPEN {
		pipeline.failure = 3
		pipeline.status = PIPELINE_STATUS_CLOSE

		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Infof("Switch status to %s", PipelineStatusNames[pipeline.status])
		return
	}

	pipeline.failure++
	if pipeline.failure >= 3 {
		pipeline.failure = 3
		pipeline.status = PIPELINE_STATUS_CLOSE

		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Infof("Switch status to %s", PipelineStatusNames[pipeline.status])
	}
}
