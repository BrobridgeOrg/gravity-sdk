package subscriber

import (
	"fmt"
	"sync"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type CollectionSnapshot struct {
	name        string
	lastKey     []byte
	isCompleted bool
}

type Snapshot struct {
	pipeline    *Pipeline
	isReady     bool
	isCompleted bool
	id          string
	collections sync.Map
}

func NewSnapshot(pipeline *Pipeline) *Snapshot {
	snapshot := &Snapshot{
		pipeline:    pipeline,
		id:          uuid.NewV1().String(),
		isCompleted: false,
	}

	// Initializing collection snapshot states
	for collectionName, _ := range pipeline.subscriber.collectionMap {
		snapshot.collections.Store(collectionName, &CollectionSnapshot{
			name:        collectionName,
			lastKey:     []byte(""),
			isCompleted: false,
		})
	}

	return snapshot
}

func (snapshot *Snapshot) getLastPosition() *CollectionSnapshot {

	var position *CollectionSnapshot = nil

	snapshot.collections.Range(func(key interface{}, value interface{}) bool {
		colSnapshot := value.(*CollectionSnapshot)

		if colSnapshot.isCompleted {
			return true
		}

		position = colSnapshot

		return false
	})

	return position
}

func (snapshot *Snapshot) UpdateLastKey(collection string, key []byte) error {
	/*
		value, ok := snapshot.collections.Load(collection)
		if !ok {
			return errors.New("No such collection")
		}

		collectionSnapshot := value.(*CollectionSnapshot)

		if bytes.Compare(collectionSnapshot.lastKey, key) == 0 {
			collectionSnapshot.isCompleted = true
		}
	*/
	//	collectionSnapshot.lastKey = key
	//TODO: Write to state store
	return nil
}

func (snapshot *Snapshot) Create() error {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.createSnapshot", snapshot.pipeline.id)
	request := pipeline_pb.CreateSnapshotRequest{
		SnapshotID: snapshot.id,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := snapshot.pipeline.subscriber.request(channel, msg, true)
	if err != nil {
		return err
	}

	var reply pipeline_pb.CreateSnapshotReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		log.Errorf("Faild to create snapshot view: %v", reply.Reason)
		return err
	}

	snapshot.isReady = true

	return nil
}

func (snapshot *Snapshot) Close() error {

	log.WithFields(logrus.Fields{
		"pipeline": snapshot.pipeline.id,
		"snapshot": snapshot.id,
	}).Info("Closing snapshot")

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.releaseSnapshot", snapshot.pipeline.id)
	request := pipeline_pb.ReleaseSnapshotRequest{
		SnapshotID: snapshot.id,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := snapshot.pipeline.subscriber.request(channel, msg, true)
	if err != nil {
		return err
	}

	var reply pipeline_pb.ReleaseSnapshotReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipeline.id,
			"snapshot": snapshot.id,
		}).Error(reply.Reason)
		return err
	}

	return nil
}

func (snapshot *Snapshot) Pull() (int64, error) {

	collectionSnapshot := snapshot.getLastPosition()
	if collectionSnapshot == nil {
		snapshot.isCompleted = true
		return 0, nil
	}

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.fetchSnapshot", snapshot.pipeline.id)
	request := pipeline_pb.FetchSnapshotRequest{
		SnapshotID:   snapshot.id,
		SubscriberID: snapshot.pipeline.subscriber.id,
		Collection:   collectionSnapshot.name,
		Key:          collectionSnapshot.lastKey,
		Offset:       1,
		Count:        int64(snapshot.pipeline.subscriber.options.ChunkSize),
	}

	if len(collectionSnapshot.lastKey) == 0 {
		request.Offset = 0
	}

	log.WithFields(logrus.Fields{
		"pipeline":   snapshot.pipeline.id,
		"snapshot":   snapshot.id,
		"collection": request.Collection,
	}).Info("Pulling from snapshot")

	msg, _ := proto.Marshal(&request)

	respData, err := snapshot.pipeline.subscriber.request(channel, msg, true)
	if err != nil {
		return 0, err
	}

	var reply pipeline_pb.FetchSnapshotReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return 0, err
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipeline.id,
			"snapshot": snapshot.id,
		}).Error(reply.Reason)
		return 0, err
	}

	// No more records for this collection
	if reply.Count == 0 {
		collectionSnapshot.isCompleted = true

		log.WithFields(logrus.Fields{
			"pipeline":   snapshot.pipeline.id,
			"snapshot":   snapshot.id,
			"collection": request.Collection,
		}).Info("No records")

		return 0, nil
	}

	log.WithFields(logrus.Fields{
		"pipeline":   snapshot.pipeline.id,
		"snapshot":   snapshot.id,
		"collection": request.Collection,
		"count":      reply.Count,
	}).Info("Receiving data from snapshot")

	collectionSnapshot.lastKey = reply.LastKey

	return reply.Count, nil
}
