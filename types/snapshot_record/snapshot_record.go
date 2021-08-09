package gravity_sdk_types_snapshot_record

import (
	"github.com/golang/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func (sr *SnapshotRecord) SetPayload(mapData map[string]interface{}) error {
	payload, err := structpb.NewStruct(mapData)
	if err != nil {
		return err
	}

	sr.Payload = payload

	return nil
}

func (sr *SnapshotRecord) ToBytes() ([]byte, error) {
	return proto.Marshal(sr)
}

func Marshal(record *SnapshotRecord) ([]byte, error) {
	return proto.Marshal(record)
}

func Unmarshal(data []byte, record *SnapshotRecord) error {
	return proto.Unmarshal(data, record)
}
