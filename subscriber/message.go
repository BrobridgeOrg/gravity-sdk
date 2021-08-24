package subscriber

import (
	"sync"

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

type MessageType int32

const (
	MESSAGE_TYPE_EVENT MessageType = iota
	MESSAGE_TYPE_SNAPSHOT
)

type Message struct {
	Pipeline     *Pipeline
	Subscription *Subscription
	Type         MessageType
	Payload      interface{}
	Callback     func(*Message)
}

func NewMessage(pipeline *Pipeline, sub *Subscription, msgType MessageType, payload interface{}) *Message {
	return &Message{
		Pipeline:     pipeline,
		Subscription: sub,
		Type:         msgType,
		Payload:      payload,
	}
}

func (msg *Message) Ack() {

	if msg.Callback != nil {
		msg.Callback(msg)
	}

	messagePool.Put(msg)
}
