package subscriber

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	ErrUnknownEventType = errors.New("pipeline: unknown event type")
)

type PipelineOpt func(*Pipeline)

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
	snapshot    *Snapshot
	id          uint64
	lastSeq     uint64
	updatedSeq  uint64
	isReady     bool
	isSuspended bool
	status      PipelineStatus
	failure     int
	mutex       sync.Mutex

	initialLoad     InitialLoadOptions
	chunkSize       int
	remoteLastSeq   uint64
	snapshotLastSeq uint64
	stateStore      StateStore
	subscriberID    string
	request         PipelineRequest
	snapshotRequest SnapshotRequest
	subscription    Subscription
	collections     []string
}

func NewPipeline(id uint64, lastSeq uint64, opts ...PipelineOpt) *Pipeline {

	p := &Pipeline{
		id:         id,
		lastSeq:    lastSeq,
		updatedSeq: lastSeq,
		isReady:    false,
		status:     PIPELINE_STATUS_OPEN,
		chunkSize:  2048,
	}

	// Initial load options by default
	p.initialLoad.Enabled = false
	p.initialLoad.Mode = "snapshot"
	p.initialLoad.OmittedCount = 100000

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithPipelineInitialLoad(enabled bool, mode string, omittedCount uint64) PipelineOpt {
	return func(p *Pipeline) {
		p.initialLoad.Enabled = enabled
		p.initialLoad.Mode = mode
		p.initialLoad.OmittedCount = omittedCount
	}
}

func WithPipelineRemoteLastSeq(lastSeq uint64) PipelineOpt {
	return func(p *Pipeline) {
		p.remoteLastSeq = lastSeq
	}
}

func WithPipelineSnapshotLastSeq(lastSeq uint64) PipelineOpt {
	return func(p *Pipeline) {
		p.snapshotLastSeq = lastSeq
	}
}

func WithPipelineStateStore(ss StateStore) PipelineOpt {
	return func(p *Pipeline) {
		p.stateStore = ss
	}
}

func WithPipelineChunkSize(size int) PipelineOpt {
	return func(p *Pipeline) {
		p.chunkSize = size
	}
}

func WithPipelineSubscriberID(id string) PipelineOpt {
	return func(p *Pipeline) {
		p.subscriberID = id
	}
}

func WithPipelineRequest(pr PipelineRequest) PipelineOpt {
	return func(p *Pipeline) {
		p.request = pr
	}
}

func WithPipelineSnapshotRequest(sr SnapshotRequest) PipelineOpt {
	return func(p *Pipeline) {
		p.snapshotRequest = sr
	}
}

func WithPipelineSubscription(s Subscription) PipelineOpt {
	return func(p *Pipeline) {
		p.subscription = s
	}
}
func WithPipelineCollections(cols []string) PipelineOpt {
	return func(p *Pipeline) {
		p.collections = cols
	}
}

func (pipeline *Pipeline) Initialize() error {

	// Initial load is disabled
	if !pipeline.initialLoad.Enabled {
		/*
			if pipeline.lastSeq == 0 {
				// Using last sequence to initialize pipeline
				pipeline.SetUpdatedSequence(pipeline.remoteLastSeq)
			}

			pipeline.lastSeq = pipeline.remoteLastSeq
		*/
		pipeline.isReady = true
		return nil
	}

	if pipeline.initialLoad.Mode == "snapshot" {

		// check states for snapshot
		if pipeline.lastSeq == 0 || pipeline.lastSeq+pipeline.initialLoad.OmittedCount < pipeline.remoteLastSeq {

			//TODO: It should be improved to compare number of records with pipeline lastS else 		//TODO: Truncate

			log.WithFields(logrus.Fields{
				"pipeline": pipeline.id,
			}).Info("Preparing snapshot")

			// Initializing snapshot
			snapshot := NewSnapshot(
				WithSnapshotPipelineID(pipeline.id),
				WithSnapshotSubscriberID(pipeline.subscriberID),
				WithSnapshotCollections(pipeline.collections),
				WithSnapshotChunkSize(pipeline.chunkSize),
				WithSnapshotRequest(pipeline.snapshotRequest),
			)
			pipeline.snapshot = snapshot

			// Using snapshot last sequence to initialize pipeline
			pipeline.lastSeq = pipeline.snapshotLastSeq

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

	if pipeline.stateStore == nil {
		return nil
	}

	// Write to store
	pipelineState, _ := pipeline.stateStore.GetPipelineState(pipeline.id)

	return pipelineState.UpdateLastSequence(updatedSeq)
}

func (pipeline *Pipeline) pull() ([][]byte, error) {

	offset := uint64(1)
	if pipeline.lastSeq == 0 {
		offset = 0
	}

	reply, err := pipeline.request.Pull(pipeline.subscriberID, pipeline.id, pipeline.lastSeq, offset, int64(pipeline.chunkSize))
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Error(err)
		return nil, err
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

	reply, err := pipeline.request.Suspend(pipeline.subscriberID, pipeline.id, pipeline.lastSeq)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Error(err)
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
	pipelineState, _ := pipeline.stateStore.GetPipelineState(pipeline.id)
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

	reply, err := pipeline.request.Awake(pipeline.subscriberID, pipeline.id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Error(err)
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

func (pipeline *Pipeline) snapshotEventCallback(msg *Message) {
	event := msg.Payload.(*SnapshotEvent)
	msg.Payload = nil
	snapshotEventPool.Put(event)
}

func (pipeline *Pipeline) eventEventCallback(msg *Message) {

	// Update sequence number to state store
	event := msg.Payload.(*DataEvent)
	msg.Payload = nil
	err := pipeline.SetUpdatedSequence(event.Sequence)
	if err != nil {
		log.Errorf("Failed to write sequence to state store: %v", err)
	}

	dataEventPool.Put(event)
}

func (pipeline *Pipeline) dispatchMessages(msgType MessageType, privData interface{}, records [][]byte) error {

	for _, record := range records {

		msg := messagePool.Get().(*Message)
		msg.Subscription = pipeline.subscription
		msg.Pipeline = pipeline
		msg.Type = msgType

		switch msg.Type {
		case MESSAGE_TYPE_SNAPSHOT:

			// Prepare data event
			snapshotEvent := snapshotEventPool.Get().(*SnapshotEvent)
			snapshotEvent.PipelineID = pipeline.id
			snapshotEvent.Collection = privData.(string)
			snapshotEvent.RawData = record
			snapshotEvent.Payload = nil

			msg.Payload = snapshotEvent
			msg.Callback = pipeline.snapshotEventCallback

		case MESSAGE_TYPE_EVENT:

			// Prepare data event
			dataEvent := dataEventPool.Get().(*DataEvent)
			dataEvent.PipelineID = pipeline.id
			dataEvent.RawData = record
			dataEvent.Sequence = 0
			dataEvent.Payload = nil

			msg.Payload = dataEvent
			msg.Callback = pipeline.eventEventCallback

		default:

			log.WithFields(logrus.Fields{
				"pipeline": pipeline.id,
			}).Warnf("pipeline: unknown event type: %d", msg.Type)

			return ErrUnknownEventType
		}

		pipeline.subscription.Push(msg)
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
