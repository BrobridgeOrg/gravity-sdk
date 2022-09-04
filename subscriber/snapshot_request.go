package subscriber

import (
	"fmt"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	"github.com/golang/protobuf/proto"
)

type SnapshotRequest interface {
	Pull(id string, sid string, pid uint64, col string, lastKey []byte, offset uint64, count int64) (*pipeline_pb.PullSnapshotReply, error)
	Create(id string, pid uint64) (*pipeline_pb.CreateSnapshotReply, error)
	Close(id string, pid uint64) (*pipeline_pb.ReleaseSnapshotReply, error)
}

func NewSnapshotRequestImpl(request RequestHandler) SnapshotRequest {
	return &SnapshotRequestImpl{
		request: request,
	}
}

type SnapshotRequestImpl struct {
	request RequestHandler
}

func (sr *SnapshotRequestImpl) Pull(id string, sid string, pid uint64, col string, lastKey []byte, offset uint64, count int64) (*pipeline_pb.PullSnapshotReply, error) {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.pullSnapshot", pid)
	request := pipeline_pb.PullSnapshotRequest{
		SnapshotID:   id,
		SubscriberID: sid,
		Collection:   col,
		Key:          lastKey,
		Offset:       offset,
		Count:        count,
	}

	msg, _ := proto.Marshal(&request)

	// Request
	respData, err := sr.request(channel, msg, true)
	if err != nil {
		return nil, fmt.Errorf("snapshot_request[Pull]: Failed to fetch: %v", err)
	}

	reply := pullSnapshotPool.Get().(*pipeline_pb.PullSnapshotReply)
	defer pullSnapshotPool.Put(reply)

	err = proto.Unmarshal(respData, reply)
	if err != nil {
		return nil, fmt.Errorf("snapshot_request[Pull]: Failed to parse response: %v", err)
	}

	return reply, nil
}

func (sr *SnapshotRequestImpl) Create(id string, pid uint64) (*pipeline_pb.CreateSnapshotReply, error) {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.createSnapshot", pid)
	request := pipeline_pb.CreateSnapshotRequest{
		SnapshotID: id,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := sr.request(channel, msg, true)
	if err != nil {
		return nil, fmt.Errorf("snapshot_request[Create]: Failed to create: %v", err)
	}

	var reply pipeline_pb.CreateSnapshotReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, fmt.Errorf("snapshot_request[Create]: Failed to parse response: %v", err)
	}

	return &reply, nil
}

func (sr *SnapshotRequestImpl) Close(id string, pid uint64) (*pipeline_pb.ReleaseSnapshotReply, error) {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.releaseSnapshot", pid)
	request := pipeline_pb.ReleaseSnapshotRequest{
		SnapshotID: id,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := sr.request(channel, msg, true)
	if err != nil {
		return nil, fmt.Errorf("snapshot_request[Close]: Failed to close: %v", err)
	}

	var reply pipeline_pb.ReleaseSnapshotReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, fmt.Errorf("snapshot_request[Close]: Failed to parse response: %v", err)
	}

	return &reply, nil
}
