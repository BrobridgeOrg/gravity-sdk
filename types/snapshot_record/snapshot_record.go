package snapshot_record

import (
	"github.com/golang/protobuf/proto"
)

func (sr *SnapshotRecord) ToBytes() ([]byte, error) {
	return proto.Marshal(sr)
}

func Marshal(record *SnapshotRecord) ([]byte, error) {
	return proto.Marshal(record)
}

func Unmarshal(data []byte, record *SnapshotRecord) error {
	return proto.Unmarshal(data, record)
}
