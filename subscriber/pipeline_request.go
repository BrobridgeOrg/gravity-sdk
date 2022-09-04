package subscriber

import (
	"fmt"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	"github.com/golang/protobuf/proto"
)

type PipelineRequest interface {
	Pull(sid string, pid uint64, startAt uint64, offset uint64, count int64) (*pipeline_pb.PullEventsReply, error)
	Suspend(sid string, pid uint64, seq uint64) (*pipeline_pb.SuspendReply, error)
	Awake(sid string, pid uint64) (*pipeline_pb.AwakeReply, error)
}

func NewPipelineRequestImpl(request RequestHandler) PipelineRequest {
	return &PipelineRequestImpl{
		request: request,
	}
}

type PipelineRequestImpl struct {
	request RequestHandler
}

func (pr *PipelineRequestImpl) Pull(sid string, pid uint64, startAt uint64, offset uint64, count int64) (*pipeline_pb.PullEventsReply, error) {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.pullEvents", pid)

	request := pipeline_pb.PullEventsRequest{
		SubscriberID: sid,
		PipelineID:   pid,
		StartAt:      startAt,
		Offset:       offset,
		Count:        count,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := pr.request(channel, msg, true)
	if err != nil {
		return nil, fmt.Errorf("pipeline_request[Pull]: Failed to fetch: %v", err)
	}

	var reply pipeline_pb.PullEventsReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, fmt.Errorf("pipeline_request[Pull]: Failed to parse response: %v", err)
	}

	return &reply, nil
}

func (pr *PipelineRequestImpl) Suspend(sid string, pid uint64, seq uint64) (*pipeline_pb.SuspendReply, error) {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.suspend", pid)

	request := pipeline_pb.SuspendRequest{
		SubscriberID: sid,
		Sequence:     seq,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := pr.request(channel, msg, true)
	if err != nil {
		return nil, fmt.Errorf("pipeline_request[Suspend]: Failed to suspend: %v", err)
	}

	var reply pipeline_pb.SuspendReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, fmt.Errorf("pipeline_request[Suspend]: Failed to parse response: %v", err)
	}

	return &reply, nil
}

func (pr *PipelineRequestImpl) Awake(sid string, pid uint64) (*pipeline_pb.AwakeReply, error) {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.awake", pid)

	request := pipeline_pb.AwakeRequest{
		SubscriberID: sid,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := pr.request(channel, msg, true)
	if err != nil {
		return nil, fmt.Errorf("pipeline_request[Awake]: Failed to awake: %v", err)
	}

	var reply pipeline_pb.AwakeReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, fmt.Errorf("pipeline_request[Awake]: Failed to response: %v", err)
	}

	return &reply, nil
}
