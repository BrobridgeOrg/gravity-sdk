package subscriber

import (
	gravity_sdk_types_event "github.com/BrobridgeOrg/gravity-sdk/types/event"
	pcf "github.com/cfsghost/parallel-chunked-flow"
	"github.com/golang/protobuf/proto"
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

			// Parsing event
			var event gravity_sdk_types_event.Event
			err := proto.Unmarshal(data.([]byte), &event)
			if err != nil {
				log.Error(err)

				// Ignore unknown event
				return
			}

			// End of chunk
			if event.Type == gravity_sdk_types_event.Event_TYPE_EVENT {
				if event.EventPayload.State == gravity_sdk_types_event.EventPayload_STATE_CHUNK_END {

					log.WithFields(logrus.Fields{
						"pipeline": event.EventPayload.PipelineID,
						"sequence": event.EventPayload.Sequence,
					}).Info("End of event chunk")

					// Pipeline chunk has no more event
					subscriber.ReleasePipeline(event.EventPayload.PipelineID)
					return
				}
			} else if event.Type == gravity_sdk_types_event.Event_TYPE_SNAPSHOT {
				if event.SnapshotInfo.State == gravity_sdk_types_event.SnapshotInfo_STATE_CHUNK_END {

					log.WithFields(logrus.Fields{
						"pipeline": event.SnapshotInfo.PipelineID,
					}).Info("End of snapshot chunck")

					// Snapshot chunk has no more event
					subscriber.ReleasePipeline(event.SnapshotInfo.PipelineID)
					return
				}
			}

			// Prepare message
			msg := messagePool.Get().(*Message)
			msg.Subscription = subscription
			msg.Event = nil
			msg.Snapshot = nil

			// Types
			switch event.Type {
			case gravity_sdk_types_event.Event_TYPE_EVENT:
				msg.Pipeline = subscriber.GetPipeline(event.EventPayload.PipelineID)
				msg.Event = event.EventPayload
			case gravity_sdk_types_event.Event_TYPE_SNAPSHOT:
				msg.Pipeline = subscriber.GetPipeline(event.SnapshotInfo.PipelineID)
				msg.Snapshot = event.SnapshotInfo
			case gravity_sdk_types_event.Event_TYPE_SYSTEM:
				msg.Pipeline = subscriber.GetPipeline(event.SystemInfo.AwakeMessage.PipelineID)

				log.WithFields(logrus.Fields{
					"pipeline": event.SystemInfo.AwakeMessage.PipelineID,
					"seq":      event.SystemInfo.AwakeMessage.Sequence,
				}).Info("Received awake event")

				subscriber.AwakePipeline(event.SystemInfo.AwakeMessage.PipelineID)
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

	if msg.Event != nil {
		s.eventHandler(msg)
	} else if msg.Snapshot != nil {
		s.snapshotHandler(msg)
	}
}

func (s *Subscription) Unsubscribe() error {
	s.buffer.Close()
	return s.sub.Unsubscribe()
}
