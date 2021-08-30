package subscriber

import (
	"errors"
	"fmt"
	"sync"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var pullSnapshotPool = sync.Pool{
	New: func() interface{} {
		return &pipeline_pb.PullSnapshotReply{}
	},
}

type CollectionSnapshot struct {
	name        string
	lastKey     []byte
	isCompleted bool
}

type Snapshot struct {
	snapshotID  string
	subscriber  *Subscriber
	pipeline    *Pipeline
	isReady     bool
	isCompleted bool
	collections sync.Map
}

func NewSnapshot(pipeline *Pipeline) *Snapshot {
	snapshot := &Snapshot{
		snapshotID: uuid.NewV1().String(),
		subscriber: pipeline.subscriber,
		pipeline:   pipeline,
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

//func (snapshot *Snapshot) pull(collection string, lastKey []byte) ([]*gravity_sdk_types_snapshot_record.SnapshotRecord, error) {
func (snapshot *Snapshot) pull(collection string, lastKey []byte) ([][]byte, error) {

	v, ok := snapshot.collections.Load(collection)
	if !ok {
		return nil, fmt.Errorf("Not subscribed to collection: %s", collection)
	}

	collectionSnapshot := v.(*CollectionSnapshot)

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.pullSnapshot", snapshot.pipeline.id)
	request := pipeline_pb.PullSnapshotRequest{
		SnapshotID:   snapshot.snapshotID,
		SubscriberID: snapshot.subscriber.id,
		Collection:   collection,
		Key:          lastKey,
		Offset:       1,
		Count:        int64(snapshot.subscriber.options.ChunkSize),
	}

	if len(lastKey) == 0 {
		request.Offset = 0
	}

	log.WithFields(logrus.Fields{
		"pipeline":   snapshot.pipeline.id,
		"snapshot":   snapshot.snapshotID,
		"collection": request.Collection,
	}).Info("Pulling from snapshot")

	msg, _ := proto.Marshal(&request)

	// Request
	respData, err := snapshot.subscriber.request(channel, msg, true)
	if err != nil {
		return nil, err
	}

	reply := pullSnapshotPool.Get().(*pipeline_pb.PullSnapshotReply)
	defer pullSnapshotPool.Put(reply)

	err = proto.Unmarshal(respData, reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipeline.id,
			"snapshot": snapshot.snapshotID,
		}).Error(reply.Reason)
		return nil, err
	}

	// No more records for this collection
	if len(reply.Records) == 0 {
		log.WithFields(logrus.Fields{
			"pipeline":   snapshot.pipeline.id,
			"snapshot":   snapshot.snapshotID,
			"collection": request.Collection,
		}).Info("Snapshot - No records")

		collectionSnapshot.isCompleted = true

		return nil, nil
	}

	// Update the last key
	collectionSnapshot.lastKey = reply.LastKey

	return reply.Records, nil
}

func (snapshot *Snapshot) Create() error {

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.createSnapshot", snapshot.pipeline.id)
	request := pipeline_pb.CreateSnapshotRequest{
		SnapshotID: snapshot.snapshotID,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := snapshot.subscriber.request(channel, msg, true)
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
		"snapshot": snapshot.snapshotID,
	}).Info("Closing snapshot")

	// Fetch events from pipelines
	channel := fmt.Sprintf("pipeline.%d.releaseSnapshot", snapshot.pipeline.id)
	request := pipeline_pb.ReleaseSnapshotRequest{
		SnapshotID: snapshot.snapshotID,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := snapshot.subscriber.request(channel, msg, true)
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
			"snapshot": snapshot.snapshotID,
		}).Error(reply.Reason)
		return err
	}

	return nil
}

//func (snapshot *Snapshot) Pull() ([]*gravity_sdk_types_snapshot_record.SnapshotRecord, error) {
func (snapshot *Snapshot) Pull() (string, [][]byte, error) {

	if !snapshot.isReady {
		return "", nil, errors.New("Snapshot is not ready yet")
	}

	// Pull snapshot chunk from one of collections
	collectionSnapshot := snapshot.getLastPosition()
	if collectionSnapshot == nil {
		snapshot.isCompleted = true
		return "", nil, nil
	}

	data, err := snapshot.pull(collectionSnapshot.name, collectionSnapshot.lastKey)

	return collectionSnapshot.name, data, err
}
