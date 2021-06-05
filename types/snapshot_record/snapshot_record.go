package gravity_sdk_types_snapshot_record

import (
	"bytes"
	"encoding/gob"
)

type SnapshotRecordMeta struct {
}

type SnapshotRecord struct {
	Meta    SnapshotRecordMeta     `json:"meta"`
	Payload map[string]interface{} `json:"payload"`
}

func (sr *SnapshotRecord) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(sr)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Unmarshal(data []byte, record *SnapshotRecord) error {
	var buf bytes.Buffer
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	return dec.Decode(record)
}
