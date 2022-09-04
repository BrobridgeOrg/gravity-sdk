package subscriber

import (
	"errors"
	"fmt"
	"strconv"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
)

func NewSnapshotRequestDummyImpl() *SnapshotRequestDummyImpl {
	return &SnapshotRequestDummyImpl{}
}

type SnapshotRequestDummyImpl struct {
	seq           uint64
	records       [][]byte
	unstableCount int
}

func (sr *SnapshotRequestDummyImpl) add(record []byte) {
	sr.seq++
	sr.records = append(sr.records, record)
}

func (sr *SnapshotRequestDummyImpl) Pull(id string, sid string, pid uint64, col string, lastKey []byte, offset uint64, count int64) (*pipeline_pb.PullSnapshotReply, error) {

	if sr.unstableCount > 0 && len(lastKey) > 0 {
		sr.unstableCount--
		return nil, errors.New("Unstable")
	}

	reply := &pipeline_pb.PullSnapshotReply{
		Success: true,
	}

	lk := uint64(0)
	if len(lastKey) > 0 {
		lk, _ = strconv.ParseUint(string(lastKey), 10, 64)
	}

	found := false
	offsetCounter := offset
	counter := count
	seq := uint64(0)
	for i, r := range sr.records {
		seq = uint64(i + 1)
		if lk == seq || lk == 0 {
			found = true
		}

		if !found {
			continue
		}

		if offsetCounter > 0 {
			offsetCounter--
			continue
		}

		reply.Records = append(reply.Records, r)
		counter--
		if counter == 0 {
			reply.LastKey = []byte(fmt.Sprintf("%d", seq))
			break
		}
	}

	return reply, nil
}

func (sr *SnapshotRequestDummyImpl) Create(id string, pid uint64) (*pipeline_pb.CreateSnapshotReply, error) {

	reply := &pipeline_pb.CreateSnapshotReply{
		Success: true,
	}

	return reply, nil
}

func (sr *SnapshotRequestDummyImpl) Close(id string, pid uint64) (*pipeline_pb.ReleaseSnapshotReply, error) {

	reply := &pipeline_pb.ReleaseSnapshotReply{
		Success: true,
	}

	return reply, nil
}
