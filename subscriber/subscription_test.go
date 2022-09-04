package subscriber

import (
	"sync"
	"testing"

	gravity_sdk_types_pipeline_event "github.com/BrobridgeOrg/gravity-sdk/types/pipeline_event"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionEventHandler(t *testing.T) {

	msgs := make(chan *Message)

	var s Subscription
	s = NewSubscriptionImpl(
		WithSubscriptionEventHandler(func(msg *Message) {
			msgs <- msg
		}),
		WithSubscriptionSnapshotHandler(func(msg *Message) {
			//event := msg.Payload.(*SnapshotEvent)
			msg.Ack()
		}),
	)

	s.Start()
	defer s.Unsubscribe()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		// Preparing original pipeline event
		event := gravity_sdk_types_pipeline_event.PipelineEvent{}
		event.PipelineID = 1
		event.Sequence = 1

		// Preparing internal message
		msg := messagePool.Get().(*Message)
		msg.Subscription = s
		msg.Pipeline = nil
		msg.Type = MESSAGE_TYPE_EVENT
		msg.Callback = func(msg *Message) {
			wg.Done()
		}

		// Preparing data event
		dataEvent := dataEventPool.Get().(*DataEvent)
		dataEvent.PipelineID = event.PipelineID
		dataEvent.Sequence = event.Sequence
		dataEvent.Payload = nil

		// Preparing payload
		record, err := gravity_sdk_types_pipeline_event.Marshal(&event)
		assert.Equal(t, nil, err)
		dataEvent.RawData = record

		msg.Payload = dataEvent

		s.Push(msg)
	}()

	for msg := range msgs {
		event := msg.Payload.(*DataEvent)
		assert.Equal(t, uint64(1), event.PipelineID)
		assert.Equal(t, uint64(1), event.Sequence)

		// Release event
		dataEventPool.Put(msg.Payload.(*DataEvent))

		msg.Payload = nil
		msg.Ack()

		break
	}

	wg.Wait()
}

func TestSubscriptionEventHandlerMultipleTimes(t *testing.T) {

	targetCount := 100000
	msgs := make(chan *Message)

	var s Subscription
	s = NewSubscriptionImpl(
		WithSubscriptionEventHandler(func(msg *Message) {
			msgs <- msg
		}),
		WithSubscriptionSnapshotHandler(func(msg *Message) {
			//event := msg.Payload.(*SnapshotEvent)
			msg.Ack()
		}),
	)

	s.Start()
	defer s.Unsubscribe()

	var wg sync.WaitGroup
	wg.Add(targetCount)

	go func() {
		for i := 0; i < targetCount; i++ {
			// Preparing original pipeline event
			event := gravity_sdk_types_pipeline_event.PipelineEvent{}
			event.PipelineID = uint64(i)
			event.Sequence = uint64(i)

			// Preparing internal message
			msg := messagePool.Get().(*Message)
			msg.Subscription = s
			msg.Pipeline = nil
			msg.Type = MESSAGE_TYPE_EVENT
			msg.Callback = func(msg *Message) {
				wg.Done()
			}

			// Preparing data event
			dataEvent := dataEventPool.Get().(*DataEvent)
			dataEvent.PipelineID = event.PipelineID
			dataEvent.Sequence = event.Sequence
			dataEvent.Payload = nil

			// Preparing payload
			record, err := gravity_sdk_types_pipeline_event.Marshal(&event)
			assert.Equal(t, nil, err)
			dataEvent.RawData = record

			msg.Payload = dataEvent

			s.Push(msg)
		}
	}()

	counter := 0
	for msg := range msgs {

		event := msg.Payload.(*DataEvent)
		assert.Equal(t, uint64(counter), event.PipelineID)
		assert.Equal(t, uint64(counter), event.Sequence)

		// Release event
		dataEventPool.Put(msg.Payload.(*DataEvent))

		msg.Payload = nil
		msg.Ack()

		counter++
		if counter == targetCount {
			break
		}
	}

	select {
	case <-msgs:
		t.Error("message is much more than expected")
	default:
	}

	wg.Wait()
}

func TestSubscriptionSnapshotHandler(t *testing.T) {

	msgs := make(chan *Message)

	var s Subscription
	s = NewSubscriptionImpl(
		WithSubscriptionEventHandler(func(msg *Message) {
			msg.Ack()
		}),
		WithSubscriptionSnapshotHandler(func(msg *Message) {
			msgs <- msg
		}),
	)

	s.Start()
	defer s.Unsubscribe()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		// Preparing original pipeline event
		event := gravity_sdk_types_pipeline_event.PipelineEvent{}
		event.PipelineID = 1

		// Preparing internal message
		msg := messagePool.Get().(*Message)
		msg.Subscription = s
		msg.Pipeline = nil
		msg.Type = MESSAGE_TYPE_SNAPSHOT
		msg.Callback = func(msg *Message) {
			wg.Done()
		}

		// Preparing snapshot event
		snapshotEvent := snapshotEventPool.Get().(*SnapshotEvent)
		snapshotEvent.PipelineID = event.PipelineID
		snapshotEvent.Collection = "test_collection"
		snapshotEvent.Payload = nil

		msg.Payload = snapshotEvent

		// Preparing payload
		record, err := gravity_sdk_types_pipeline_event.Marshal(&event)
		assert.Equal(t, nil, err)

		snapshotEvent.RawData = record

		msg.Payload = snapshotEvent

		s.Push(msg)
	}()

	for msg := range msgs {
		event := msg.Payload.(*SnapshotEvent)
		assert.Equal(t, uint64(1), event.PipelineID)
		assert.Equal(t, "test_collection", event.Collection)

		// Release event
		snapshotEventPool.Put(msg.Payload.(*SnapshotEvent))

		msg.Payload = nil
		msg.Ack()

		break
	}

	wg.Wait()
}

func TestSubscriptionSnapshotHandlerMultipleTimes(t *testing.T) {

	targetCount := 100000
	msgs := make(chan *Message)

	var s Subscription
	s = NewSubscriptionImpl(
		WithSubscriptionEventHandler(func(msg *Message) {
			msg.Ack()
		}),
		WithSubscriptionSnapshotHandler(func(msg *Message) {
			msgs <- msg
		}),
	)

	s.Start()
	defer s.Unsubscribe()

	var wg sync.WaitGroup
	wg.Add(targetCount)

	go func() {
		for i := 0; i < targetCount; i++ {
			// Preparing original pipeline event
			event := gravity_sdk_types_pipeline_event.PipelineEvent{}
			event.PipelineID = uint64(i)

			// Preparing internal message
			msg := messagePool.Get().(*Message)
			msg.Subscription = s
			msg.Pipeline = nil
			msg.Type = MESSAGE_TYPE_SNAPSHOT
			msg.Callback = func(msg *Message) {
				wg.Done()
			}

			// Preparing snapshot event
			snapshotEvent := snapshotEventPool.Get().(*SnapshotEvent)
			snapshotEvent.PipelineID = event.PipelineID
			snapshotEvent.Collection = "test_collection"
			snapshotEvent.Payload = nil

			msg.Payload = snapshotEvent

			// Preparing payload
			record, err := gravity_sdk_types_pipeline_event.Marshal(&event)
			assert.Equal(t, nil, err)

			snapshotEvent.RawData = record

			msg.Payload = snapshotEvent

			s.Push(msg)
		}
	}()

	counter := 0
	for msg := range msgs {
		event := msg.Payload.(*SnapshotEvent)
		assert.Equal(t, uint64(counter), event.PipelineID)
		assert.Equal(t, "test_collection", event.Collection)

		// Release event
		snapshotEventPool.Put(msg.Payload.(*SnapshotEvent))

		msg.Payload = nil
		msg.Ack()

		counter++
		if counter == targetCount {
			break
		}
	}

	wg.Wait()
}
