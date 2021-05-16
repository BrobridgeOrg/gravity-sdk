package subscriber

import (
	gravity_sdk_types_event "github.com/BrobridgeOrg/gravity-sdk/types/event"
	pcf "github.com/cfsghost/parallel-chunked-flow"
	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
)

type Subscription struct {
	subscriber *Subscriber
	sub        *nats.Subscription
	buffer     *pcf.ParallelChunkedFlow
	callback   MessageHandler
}

func NewSubscription(subscriber *Subscriber, bufferSize int) *Subscription {

	subscription := &Subscription{
		subscriber: subscriber,
	}

	// Create Options object
	options := &pcf.Options{
		BufferSize: bufferSize,
		ChunkSize:  1024,
		ChunkCount: 128,
		Handler: func(data interface{}, output chan interface{}) {

			// Parsing event
			var event gravity_sdk_types_event.Event
			err := proto.Unmarshal(data.([]byte), &event)
			if err != nil {
				log.Error(err)

				// Ignore unknown event
				return
			}

			output <- &event
		},
	}

	subscription.buffer = pcf.NewParallelChunkedFlow(options)

	return subscription
}

func (s *Subscription) start() {
	go func() {
		for event := range s.buffer.Output() {
			s.callback(event.(*gravity_sdk_types_event.Event))
		}
	}()
}

func (s *Subscription) Unsubscribe() error {
	s.buffer.Close()
	return s.sub.Unsubscribe()
}
