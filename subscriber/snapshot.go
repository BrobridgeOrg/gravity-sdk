package subscriber

import (
	"errors"
	"fmt"
	"sync"

	pipeline_pb "github.com/BrobridgeOrg/gravity-api/service/pipeline"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var pullSnapshotPool = sync.Pool{
	New: func() interface{} {
		return &pipeline_pb.PullSnapshotReply{}
	},
}

type SnapshotOpt func(*Snapshot)

type CollectionSnapshot struct {
	name        string
	lastKey     []byte
	isCompleted bool
}

type Snapshot struct {
	snapshotID   string
	pipelineID   uint64
	subscriberID string
	isReady      bool
	isCompleted  bool
	collections  sync.Map
	request      SnapshotRequest
	chunkSize    int
}

func NewSnapshot(opts ...SnapshotOpt) *Snapshot {
	snapshot := &Snapshot{
		snapshotID: uuid.NewV1().String(),
	}

	for _, opt := range opts {
		opt(snapshot)
	}

	return snapshot
}

func WithSnapshotPipelineID(id uint64) SnapshotOpt {
	return func(s *Snapshot) {
		s.pipelineID = id
	}
}

func WithSnapshotSubscriberID(id string) SnapshotOpt {
	return func(s *Snapshot) {
		s.subscriberID = id
	}
}

func WithSnapshotCollections(cols []string) SnapshotOpt {
	return func(s *Snapshot) {

		// Initializing collection snapshot states
		for _, collectionName := range cols {
			s.collections.Store(collectionName, &CollectionSnapshot{
				name:        collectionName,
				lastKey:     []byte(""),
				isCompleted: false,
			})
		}
	}
}

func WithSnapshotRequest(sr SnapshotRequest) SnapshotOpt {
	return func(s *Snapshot) {
		s.request = sr
	}
}

func WithSnapshotChunkSize(size int) SnapshotOpt {
	return func(s *Snapshot) {
		s.chunkSize = size
	}
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

	log.WithFields(logrus.Fields{
		"pipeline":   snapshot.pipelineID,
		"snapshot":   snapshot.subscriberID,
		"collection": collection,
	}).Info("Pulling from snapshot")

	offset := uint64(1)
	if len(lastKey) == 0 {
		offset = 0
	}

	reply, err := snapshot.request.Pull(
		snapshot.snapshotID,
		snapshot.subscriberID,
		snapshot.pipelineID,
		collection,
		lastKey,
		offset,
		int64(snapshot.chunkSize),
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipelineID,
			"snapshot": snapshot.snapshotID,
		}).Error(err)
		return nil, err
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipelineID,
			"snapshot": snapshot.snapshotID,
		}).Error(reply.Reason)
		return nil, err
	}

	// No more records for this collection
	if len(reply.Records) == 0 {
		log.WithFields(logrus.Fields{
			"pipeline":   snapshot.pipelineID,
			"snapshot":   snapshot.snapshotID,
			"collection": collection,
		}).Info("Snapshot - No records")

		collectionSnapshot.isCompleted = true

		return nil, nil
	}

	// Update the last key
	collectionSnapshot.lastKey = reply.LastKey

	return reply.Records, nil
}

func (snapshot *Snapshot) Create() error {

	reply, err := snapshot.request.Create(snapshot.snapshotID, snapshot.pipelineID)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipelineID,
			"snapshot": snapshot.snapshotID,
		}).Error(err)
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
		"pipeline": snapshot.pipelineID,
		"snapshot": snapshot.snapshotID,
	}).Info("Closing snapshot")

	reply, err := snapshot.request.Close(snapshot.snapshotID, snapshot.pipelineID)
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipelineID,
			"snapshot": snapshot.snapshotID,
		}).Error(err)
		return err
	}

	if !reply.Success {
		log.WithFields(logrus.Fields{
			"pipeline": snapshot.pipelineID,
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
