package subscriber

import (
	"errors"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
)

func NewPipelineRequestDummyImpl() *PipelineRequestDummyImpl {
	return &PipelineRequestDummyImpl{
		events: make([][]byte, 0),
	}
}

type PipelineRequestDummyImpl struct {
	seq           uint64
	events        [][]byte
	unstableCount int
}

func (pr *PipelineRequestDummyImpl) add(event []byte) {
	pr.seq++
	pr.events = append(pr.events, event)
}

func (pr *PipelineRequestDummyImpl) Pull(sid string, pid uint64, startAt uint64, offset uint64, count int64) (*pipeline_pb.PullEventsReply, error) {

	if pr.unstableCount > 0 && startAt > 0 {
		pr.unstableCount--
		return nil, errors.New("Unstable")
	}

	reply := &pipeline_pb.PullEventsReply{
		Success: true,
	}

	found := false
	offsetCounter := offset
	counter := count
	seq := uint64(0)
	for i, event := range pr.events {
		seq = uint64(i + 1)

		if startAt == seq || startAt == 0 {
			found = true
		}

		if !found {
			continue
		}

		if offsetCounter > 0 {
			offsetCounter--
			continue
		}

		reply.Events = append(reply.Events, event)
		counter--
		if counter == 0 {
			reply.LastSeq = seq
			break
		}
	}

	return reply, nil
}

func (pr *PipelineRequestDummyImpl) Suspend(sid string, pid uint64, seq uint64) (*pipeline_pb.SuspendReply, error) {

	reply := &pipeline_pb.SuspendReply{
		Success: true,
	}

	return reply, nil
}

func (pr *PipelineRequestDummyImpl) Awake(sid string, pid uint64) (*pipeline_pb.AwakeReply, error) {

	reply := &pipeline_pb.AwakeReply{
		Success: true,
	}

	return reply, nil
}
