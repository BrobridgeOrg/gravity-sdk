package subscriber

import (
	"sync"

	gravity_sdk_types_event "github.com/BrobridgeOrg/gravity-sdk/types/event"
	gravity_sdk_types_projection "github.com/BrobridgeOrg/gravity-sdk/types/projection"
)

var projectionPool = sync.Pool{
	New: func() interface{} {
		return &gravity_sdk_types_projection.Projection{}
	},
}

var messagePool = sync.Pool{
	New: func() interface{} {
		return &Message{}
	},
}

type Message struct {
	Pipeline     *Pipeline
	Subscription *Subscription
	Event        *gravity_sdk_types_event.EventPayload
	Snapshot     *gravity_sdk_types_event.SnapshotInfo
}

func (msg *Message) Ack() {

	if msg.Event != nil {

		// Update state store
		if msg.Subscription.subscriber.options.StateStore != nil {
			pipelineState, _ := msg.Subscription.subscriber.options.StateStore.GetPipelineState(msg.Event.PipelineID)
			pipelineState.UpdateLastSequence(msg.Event.Sequence)
		}
	} else if msg.Snapshot != nil {
		if msg.Snapshot.State == gravity_sdk_types_event.SnapshotInfo_STATE_NONE {

			// update last key
			msg.Pipeline.snapshot.UpdateLastKey(msg.Snapshot.Collection, msg.Snapshot.Key)
		}
	}

	messagePool.Put(msg)
}
