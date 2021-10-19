package subscriber

import (
	"sync"

	gravity_sdk_types_record "github.com/BrobridgeOrg/gravity-sdk/types/record"
	gravity_sdk_types_snapshot_record "github.com/BrobridgeOrg/gravity-sdk/types/snapshot_record"
)

type DataEvent struct {
	PipelineID uint64
	Sequence   uint64
	RawData    []byte
	//	Payload    *gravity_sdk_types_projection.Projection
	Payload *gravity_sdk_types_record.Record
}

type SnapshotEvent struct {
	PipelineID uint64
	Collection string
	RawData    []byte
	Payload    *gravity_sdk_types_snapshot_record.SnapshotRecord
}

var dataEventPool = sync.Pool{
	New: func() interface{} {
		return &DataEvent{}
	},
}

var snapshotEventPool = sync.Pool{
	New: func() interface{} {
		return &SnapshotEvent{}
	},
}
