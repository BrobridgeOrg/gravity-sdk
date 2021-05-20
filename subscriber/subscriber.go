package subscriber

import (
	"fmt"
	"os"
	"time"

	subscriber_manager_pb "github.com/BrobridgeOrg/gravity-api/service/subscriber_manager"
	synchronizer_pb "github.com/BrobridgeOrg/gravity-api/service/synchronizer"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/BrobridgeOrg/gravity-sdk/pipeline_manager"
	"github.com/BrobridgeOrg/gravity-sdk/subscriber_manager"
	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const (
	SubscriberType_Transmitter subscriber_manager_pb.SubscriberType = subscriber_manager_pb.SubscriberType_TRANSMITTER
	SubscriberType_Exporter    subscriber_manager_pb.SubscriberType = subscriber_manager_pb.SubscriberType_EXPORTER
)

type MessageHandler func(*Message)

type Subscriber struct {
	client        *core.Client
	host          string
	options       *Options
	id            string
	pipelines     []*Pipeline
	collectionMap map[string][]string

	tasks chan *Pipeline
}

func NewSubscriber(options *Options) *Subscriber {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	return &Subscriber{
		options:       options,
		pipelines:     make([]*Pipeline, 0),
		collectionMap: make(map[string][]string),
	}
}

func NewSubscriberWithClient(client *core.Client, options *Options) *Subscriber {

	subscriber := NewSubscriber(options)
	subscriber.client = client

	return subscriber
}

func (sub *Subscriber) register(subscriberType subscriber_manager_pb.SubscriberType, component string, subscriberID string, name string) error {

	log.WithFields(logrus.Fields{
		"id": subscriberID,
	}).Info("Registering subscriber")

	conn := sub.client.GetConnection()

	request := subscriber_manager_pb.RegisterSubscriberRequest{
		SubscriberID: subscriberID,
		Name:         name,
		Type:         subscriberType,
		Component:    component,
	}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request("gravity.subscriber_manager.registerSubscriber", msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply subscriber_manager_pb.RegisterSubscriberReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	sub.id = subscriberID

	log.WithFields(logrus.Fields{
		"id": subscriberID,
	}).Info("Subscriber was registered")

	return nil
}

func (sub *Subscriber) fetch(pipelineID uint64, startAt uint64) (uint64, uint64, error) {

	conn := sub.client.GetConnection()

	// Fetch events from pipelines
	channel := fmt.Sprintf("gravity.pipeline.%d.fetch", pipelineID)
	/*
		log.WithFields(logrus.Fields{
			"pipeline": pipelineID,
		}).Info("Fetching data from pipeline")
	*/
	request := synchronizer_pb.PipelineFetchRequest{
		SubscriberID: sub.id,
		PipelineID:   pipelineID,
		StartAt:      startAt,
		Offset:       1,
		Count:        int64(sub.options.ChunkSize),
	}

	if startAt == 0 {
		request.Offset = 0
	}

	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request(channel, msg, time.Second*10)
	if err != nil {
		return 0, startAt, err
	}

	var reply synchronizer_pb.PipelineFetchReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return 0, startAt, err
	}

	if !reply.Success {
		log.Error(reply.Reason)
		return 0, startAt, err
	}

	return reply.Count, reply.LastSeq, nil
}

func (sub *Subscriber) startWorker(workerID int) {

	log.WithFields(logrus.Fields{
		"worker": workerID,
	}).Info("Started wroker")

	for {
		select {
		case pipeline := <-sub.tasks:

			// Fetching data from specfic pipeline
			count, lastSeq, err := sub.fetch(pipeline.id, pipeline.lastSeq)
			if err != nil {
				log.WithFields(logrus.Fields{
					"worker":   workerID,
					"pipeline": pipeline.id,
				}).Error(err)
				continue
			}

			// Update pipeline state
			pipeline.lastSeq = lastSeq

			if count > 0 {
				log.WithFields(logrus.Fields{
					"worker":   workerID,
					"pipeline": pipeline.id,
					"count":    count,
				}).Info("Received records")
			}

			// re-queue
			sub.tasks <- pipeline
		}
	}
}

func (sub *Subscriber) prepareWorkers() {

	for i := 0; i < sub.options.WorkerCount; i++ {
		go sub.startWorker(i)
	}

	sub.tasks = make(chan *Pipeline, len(sub.pipelines))

	for _, pipeline := range sub.pipelines {
		sub.tasks <- pipeline
	}
}

func (sub *Subscriber) Connect(host string, options *core.Options) error {
	sub.client = core.NewClient()
	return sub.client.Connect(host, options)
}

func (sub *Subscriber) Disconnect() {
	sub.client.Disconnect()
}

func (sub *Subscriber) Register(subscriberType subscriber_manager_pb.SubscriberType, component string, subscriberID string, name string) error {

	id := subscriberID
	if len(id) == 0 {
		id = uuid.NewV1().String()
	}

	// Register to get channel id
	err := sub.register(subscriberType, component, id, name)
	if err != nil {
		return err
	}

	return nil
}

func (sub *Subscriber) Subscribe(cb MessageHandler) (*Subscription, error) {

	channel := fmt.Sprintf("gravity.subscriber.%s", sub.id)

	log.WithFields(logrus.Fields{
		"channel": channel,
	}).Info("Subscribe to synchonizer")

	subscription := NewSubscription(sub, sub.options.BufferSize)
	subscription.callback = cb

	// Subscribe to channel
	conn := sub.client.GetConnection()
	s, err := conn.Subscribe(channel, func(m *nats.Msg) {
		subscription.buffer.Push(m.Data)
	})
	if err != nil {
		return nil, err
	}

	subscription.sub = s
	subscription.start()

	// Start to receive data
	//go sub.poll(count)
	sub.prepareWorkers()

	return subscription, nil
}

func (sub *Subscriber) GetPipelineCount() (uint64, error) {

	// Getting pipeline count
	pm := pipeline_manager.NewPipelineManagerWithClient(sub.client, pipeline_manager.NewOptions())
	return pm.GetPipelineCount()
}

func (sub *Subscriber) AddAllPipelines() error {

	// Getting pipeline count
	count, err := sub.GetPipelineCount()
	if err != nil {
		return err
	}

	for i := uint64(0); i < count; i++ {
		pipeline := NewPipeline(i, 0)
		err := sub.AddPipeline(pipeline)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sub *Subscriber) AddPipeline(pipeline *Pipeline) error {

	pipelineState, err := sub.options.StateStore.GetPipelineState(pipeline.id)
	if err != nil {
		return err
	}

	// Load state
	log.WithFields(logrus.Fields{
		"pipeline": pipeline.id,
		"lastSeq":  pipelineState.GetLastSequence(),
	}).Info("Loaded pipeline state")
	pipeline.UpdateLastSequence(pipelineState.GetLastSequence())

	sub.pipelines = append(sub.pipelines, pipeline)
	return nil
}

func (sub *Subscriber) SubscribeToCollections(colMap map[string][]string) error {

	if len(colMap) == 0 {
		return nil
	}

	// Subscribe to collections
	collections := make([]string, 0, len(colMap))
	for collectionName, tables := range colMap {
		sub.collectionMap[collectionName] = tables
		collections = append(collections, collectionName)
	}

	// Call controller to subscribe
	sm := subscriber_manager.NewSubscriberManagerWithClient(sub.client, subscriber_manager.NewOptions())
	return sm.SubscribeToCollections(sub.id, collections)
}

func (sub *Subscriber) GetCollectionInfo(collection string) []string {
	tables, ok := sub.collectionMap[collection]
	if !ok {
		return nil
	}

	return tables
}
