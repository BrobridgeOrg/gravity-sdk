package subscriber

import (
	"fmt"
	"time"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	synchronizer_pb "github.com/BrobridgeOrg/gravity-api/service/synchronizer"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

type Pipeline struct {
	subscriber    *Subscriber
	snapshot      *Snapshot
	id            uint64
	lastSeq       uint64
	isInitialized bool
	isReady       bool
}

func NewPipeline(subscriber *Subscriber, id uint64, lastSeq uint64) *Pipeline {
	return &Pipeline{
		subscriber:    subscriber,
		id:            id,
		lastSeq:       lastSeq,
		isInitialized: false,
		isReady:       false,
	}
}

func (pipeline *Pipeline) UpdateLastSequence(sequence uint64) {
	pipeline.lastSeq = sequence
}

func (pipeline *Pipeline) initializeState() error {

	conn := pipeline.subscriber.client.GetConnection()

	// Fetch events from pipelines
	channel := fmt.Sprintf("gravity.pipeline.%d.getState", pipeline.id)
	request := pipeline_pb.GetStateRequest{}

	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request(channel, msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply pipeline_pb.GetStateReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		log.Error(reply.Reason)
		return err
	}

	pipeline.lastSeq = reply.LastSeq
	log.Warn(pipeline.lastSeq)

	return nil

}

func (pipeline *Pipeline) performInitialLoad() error {

	if pipeline.snapshot == nil {

		// Getting pipeline states
		err := pipeline.initializeState()
		if err != nil {
			return err
		}

		// TODO: check states
		//pipeline.subscriber.options.InitialLoad.OmittedCount

		// Create a new snapshot
		snapshot := NewSnapshot(pipeline)
		pipeline.snapshot = snapshot
		err = snapshot.Create()
		if err != nil {
			return err
		}
	}

	// Pull data
	pipeline.snapshot.Pull()

	return nil
}

func (pipeline *Pipeline) fetch() error {

	conn := pipeline.subscriber.client.GetConnection()

	// Fetch events from pipelines
	channel := fmt.Sprintf("gravity.pipeline.%d.fetch", pipeline.id)
	/*
		log.WithFields(logrus.Fields{
			"pipeline": pipelineID,
		}).Info("Fetching data from pipeline")
	*/
	request := synchronizer_pb.PipelineFetchRequest{
		SubscriberID: pipeline.subscriber.id,
		PipelineID:   pipeline.id,
		StartAt:      pipeline.lastSeq,
		Offset:       1,
		Count:        int64(pipeline.subscriber.options.ChunkSize),
	}

	if pipeline.lastSeq == 0 {
		request.Offset = 0
	}

	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request(channel, msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply synchronizer_pb.PipelineFetchReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		log.Error(reply.Reason)
		return err
	}

	pipeline.lastSeq = reply.LastSeq

	if reply.Count > 0 {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
			"count":    reply.Count,
		}).Info("Received records")
	}

	return nil
}

func (pipeline *Pipeline) Pull() error {

	if !pipeline.isInitialized {

		// Initial load is enabled
		if pipeline.subscriber.options.InitialLoad.Enabled {
			err := pipeline.performInitialLoad()
			if err != nil {
				return err
			}
		}

		pipeline.isInitialized = true
		//pipeline.isReady = true
	}

	if pipeline.isReady {
		return pipeline.fetch()
	}

	return nil
}
