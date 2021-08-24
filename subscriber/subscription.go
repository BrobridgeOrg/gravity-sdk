package subscriber

import (
	gravity_sdk_types_pipeline_event "github.com/BrobridgeOrg/gravity-sdk/types/pipeline_event"
	gravity_sdk_types_projection "github.com/BrobridgeOrg/gravity-sdk/types/projection"
	gravity_sdk_types_snapshot_record "github.com/BrobridgeOrg/gravity-sdk/types/snapshot_record"
	pcf "github.com/cfsghost/parallel-chunked-flow"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type Subscription struct {
	subscriber      *Subscriber
	sub             *nats.Subscription
	buffer          *pcf.ParallelChunkedFlow
	eventHandler    MessageHandler
	snapshotHandler MessageHandler
}

func NewSubscription(subscriber *Subscriber, bufferSize int) *Subscription {

	subscription := &Subscription{
		subscriber: subscriber,
		eventHandler: func(msg *Message) {
			msg.Ack()
		},
		snapshotHandler: func(msg *Message) {
			msg.Ack()
		},
	}

	// Create Options object
	options := &pcf.Options{
		BufferSize: bufferSize,
		ChunkSize:  1024,
		ChunkCount: 128,
		Handler: func(data interface{}, done func(interface{})) {

			msg := data.(*Message)

			switch msg.Type {
			case MESSAGE_TYPE_SNAPSHOT:

				event := msg.Payload.(*SnapshotEvent)

				// Parsing snapshot record
				var snapshotRecord gravity_sdk_types_snapshot_record.SnapshotRecord
				err := gravity_sdk_types_snapshot_record.Unmarshal(event.RawData, &snapshotRecord)
				if err != nil {
					log.Error(err)
					return
				}

				event.Payload = &snapshotRecord
			case MESSAGE_TYPE_EVENT:

				event := msg.Payload.(*DataEvent)

				// Parsing event
				var pe gravity_sdk_types_pipeline_event.PipelineEvent
				err := gravity_sdk_types_pipeline_event.Unmarshal(event.RawData, &pe)
				if err != nil {
					log.WithFields(logrus.Fields{
						"pipeline": event.PipelineID,
					}).Errorf("pipeline event - %v", err)
					return
				}

				var pj gravity_sdk_types_projection.Projection
				err = gravity_sdk_types_projection.Unmarshal(pe.Payload, &pj)
				if err != nil {
					log.WithFields(logrus.Fields{
						"pipeline": event.PipelineID,
					}).Error(err)
					return
				}

				event.Sequence = pe.Sequence
				event.Payload = &pj
			}

			done(msg)
		},
	}

	subscription.buffer = pcf.NewParallelChunkedFlow(options)

	return subscription
}

func (s *Subscription) start() {
	go func() {
		for msg := range s.buffer.Output() {
			s.handle(msg.(*Message))
		}
	}()
}

func (s *Subscription) handle(msg *Message) {

	switch msg.Type {
	case MESSAGE_TYPE_EVENT:
		s.eventHandler(msg)
	case MESSAGE_TYPE_SNAPSHOT:
		s.snapshotHandler(msg)
	}
}

func (s *Subscription) Push(msg *Message) {
	s.buffer.Push(msg)
}

func (s *Subscription) Unsubscribe() error {
	s.buffer.Close()
	return s.sub.Unsubscribe()
}
