package subscriber

import (
	gravity_sdk_types_projection "github.com/BrobridgeOrg/gravity-sdk/types/projection"
	gravity_sdk_types_snapshot_record "github.com/BrobridgeOrg/gravity-sdk/types/snapshot_record"
)

type DataEvent struct {
	PipelineID uint64
	Sequence   uint64
	RawData    []byte
	Payload    *gravity_sdk_types_projection.Projection
}

type SnapshotEvent struct {
	PipelineID uint64
	Collection string
	RawData    []byte
	Payload    *gravity_sdk_types_snapshot_record.SnapshotRecord
}
