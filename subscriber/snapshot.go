package subscriber

import (
	"fmt"
	"time"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

type Snapshot struct {
	pipeline *Pipeline
	isReady  bool
	id       string
	count    uint64
	total    uint64
	lastKey  string
}

func NewSnapshot(pipeline *Pipeline) *Snapshot {
	return &Snapshot{
		pipeline: pipeline,
		id:       uuid.NewV1().String(),
	}
}

func (snapshot *Snapshot) Create() error {

	conn := snapshot.pipeline.subscriber.client.GetConnection()

	// Fetch events from pipelines
	channel := fmt.Sprintf("gravity.pipeline.%d.createSnapshot", snapshot.pipeline.id)
	request := pipeline_pb.CreateSnapshotRequest{
		SnapshotID: snapshot.id,
	}

	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request(channel, msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply pipeline_pb.CreateSnapshotReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		log.Error(reply.Reason)
		return err
	}

	snapshot.isReady = true

	return nil
}

func (snapshot *Snapshot) Pull() error {
	/*
		conn := snapshot.pipeline.subscriber.client.GetConnection()

		// Fetch events from pipelines
		channel := fmt.Sprintf("gravity.pipeline.%d.fetchSnapshot", snapshot.pipeline.id)
		request := pipeline_pb.FetchSnapshotRequest{
			SnapshotID: snapshot.id,
			Key:        snapshot.lastKey,
			Offset:     0,
		}

		msg, _ := proto.Marshal(&request)

		resp, err := conn.Request(channel, msg, time.Second*10)
		if err != nil {
			return err
		}

		var reply pipeline_pb.FetchSnapshotReply
		err = proto.Unmarshal(resp.Data, &reply)
		if err != nil {
			return err
		}

		if !reply.Success {
			log.Error(reply.Reason)
			return err
		}
	*/
	return nil
}
