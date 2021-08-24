package gravity_sdk_types_event

import (
	"github.com/golang/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func (event *Event) SetPayload(mapData map[string]interface{}) error {
	payload, err := structpb.NewStruct(mapData)
	if err != nil {
		return err
	}

	event.Payload = payload

	return nil
}

func Marshal(event *Event) ([]byte, error) {
	return proto.Marshal(event)
}

func Unmarshal(data []byte, event *Event) error {
	return proto.Unmarshal(data, event)
}
