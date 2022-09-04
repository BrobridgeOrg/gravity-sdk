package subscriber

import (
	"sync"

	gravity_sdk_types_pipeline_event "github.com/BrobridgeOrg/gravity-sdk/types/pipeline_event"
	gravity_sdk_types_record "github.com/BrobridgeOrg/gravity-sdk/types/record"

	gravity_sdk_types_snapshot_record "github.com/BrobridgeOrg/gravity-sdk/types/snapshot_record"
	sdf "github.com/BrobridgeOrg/sequential-data-flow"
	"github.com/sirupsen/logrus"
)

var pipelineEventPool = sync.Pool{
	New: func() interface{} {
		return &gravity_sdk_types_pipeline_event.PipelineEvent{}
	},
}

var snapshotRecordPool = sync.Pool{
	New: func() interface{} {
		return &gravity_sdk_types_snapshot_record.SnapshotRecord{}
	},
}

var recordPool = sync.Pool{
	New: func() interface{} {
		return &gravity_sdk_types_record.Record{}
	},
}

type SubscriptionOpt func(s *SubscriptionImpl)

type Subscription interface {
	Start()
	Push(msg *Message)
	Unsubscribe() error
}

type SubscriptionImpl struct {
	buffer          *sdf.Flow
	bufferSize      int
	workerCount     int
	eventHandler    MessageHandler
	snapshotHandler MessageHandler
}

func NewSubscriptionImpl(opts ...SubscriptionOpt) *SubscriptionImpl {

	subscription := &SubscriptionImpl{
		bufferSize:  10240,
		workerCount: 4,
	}

	for _, opt := range opts {
		opt(subscription)
	}

	// Initializing sequential data flow
	options := sdf.NewOptions()
	options.BufferSize = 10240
	options.WorkerCount = subscription.workerCount
	options.Handler = subscription.prepare

	subscription.buffer = sdf.NewFlow(options)

	return subscription
}

func WithSubscriptionBufferSize(size int) SubscriptionOpt {
	return func(s *SubscriptionImpl) {
		s.bufferSize = size
	}
}

func WithSubscriptionWorkerCount(count int) SubscriptionOpt {
	return func(s *SubscriptionImpl) {
		s.workerCount = count
	}
}

func WithSubscriptionEventHandler(fn MessageHandler) SubscriptionOpt {
	return func(s *SubscriptionImpl) {
		s.eventHandler = fn
	}
}

func WithSubscriptionSnapshotHandler(fn MessageHandler) SubscriptionOpt {
	return func(s *SubscriptionImpl) {
		s.snapshotHandler = fn
	}
}

func (s *SubscriptionImpl) Start() {
	go func() {
		for msg := range s.buffer.Output() {

			// Skip
			if msg == nil {
				continue
			}

			s.handle(msg.(*Message))
		}
	}()
}

func (s *SubscriptionImpl) prepare(data interface{}, done func(interface{})) {

	msg := data.(*Message)

	switch msg.Type {
	case MESSAGE_TYPE_SNAPSHOT:

		event := msg.Payload.(*SnapshotEvent)

		// Parsing snapshot record
		snapshotRecord := snapshotRecordPool.Get().(*gravity_sdk_types_snapshot_record.SnapshotRecord)
		err := gravity_sdk_types_snapshot_record.Unmarshal(event.RawData, snapshotRecord)
		if err != nil {
			log.Error(err)
			done(nil)
			return
		}

		event.Payload = snapshotRecord
	case MESSAGE_TYPE_EVENT:

		event := msg.Payload.(*DataEvent)

		// Parsing event
		pe := pipelineEventPool.Get().(*gravity_sdk_types_pipeline_event.PipelineEvent)
		err := gravity_sdk_types_pipeline_event.Unmarshal(event.RawData, pe)
		if err != nil {
			log.WithFields(logrus.Fields{
				"pipeline": event.PipelineID,
			}).Errorf("pipeline event - %v", err)
			done(nil)
			return
		}

		record := recordPool.Get().(*gravity_sdk_types_record.Record)
		err = gravity_sdk_types_record.Unmarshal(pe.Payload, record)
		if err != nil {
			log.WithFields(logrus.Fields{
				"pipeline": event.PipelineID,
			}).Error(err)
			done(nil)
			return
		}

		event.Sequence = pe.Sequence
		event.Payload = record

		pipelineEventPool.Put(pe)
	default:
		log.Warnf("subscription: unknown event type: %d", msg.Type)
		done(nil)
		return
	}

	done(msg)
}

func (s *SubscriptionImpl) handle(msg *Message) {

	switch msg.Type {
	case MESSAGE_TYPE_EVENT:
		s.eventHandler(msg)
	case MESSAGE_TYPE_SNAPSHOT:
		s.snapshotHandler(msg)
	}
}

func (s *SubscriptionImpl) Push(msg *Message) {
	s.buffer.Push(msg)
}

func (s *SubscriptionImpl) Unsubscribe() error {
	s.buffer.Close()
	return nil
}
