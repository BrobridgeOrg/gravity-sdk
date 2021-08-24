package subscriber

import (
	"errors"

	gravity_sdk_types_event "github.com/BrobridgeOrg/gravity-sdk/types/event"
)

type EventHandler struct {
	subscriber *Subscriber
}

func NewEventHandler(subscriber *Subscriber) *EventHandler {
	return &EventHandler{
		subscriber: subscriber,
	}
}

func (eh *EventHandler) ProcessEvent(data []byte) error {

	var event gravity_sdk_types_event.Event
	err := gravity_sdk_types_event.Unmarshal(data, &event)
	if err != nil {
		return err
	}

	switch event.EventName {
	case "awake":
		payload := event.Payload.AsMap()

		var id uint64
		v, ok := payload["pipelineID"]
		if !ok {
			return errors.New("Invalid event")
		}

		id = uint64(v.(float64))

		//		var lastSeq uint64
		_, ok = payload["lastSeq"]
		if !ok {
			return errors.New("Invalid event")
		}

		//		lastSeq = uint64(v.(float64))

		//eh.subscriber.scheduler.Awake(id, lastSeq)
		eh.subscriber.runner.Awake(id)
	}

	return nil
}
