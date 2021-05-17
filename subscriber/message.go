package subscriber

import (
	"sync"

	gravity_sdk_types_event "github.com/BrobridgeOrg/gravity-sdk/types/event"
)

var messagePool = sync.Pool{
	New: func() interface{} {
		return &Message{}
	},
}

type Message struct {
	PipelineID   uint64
	Subscription *Subscription
	Event        *gravity_sdk_types_event.Event
}

func (msg *Message) Ack() {

	// Update state store
	if msg.Subscription.subscriber.options.StateStore != nil {
		pipelineState, _ := msg.Subscription.subscriber.options.StateStore.GetPipelineState(msg.PipelineID)
		pipelineState.UpdateLastSequence(msg.Event.Sequence)
	}

	messagePool.Put(msg)
}
