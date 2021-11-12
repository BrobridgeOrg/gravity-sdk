package subscriber

import (
	"fmt"
	"os"
	"time"

	subscriber_manager_pb "github.com/BrobridgeOrg/gravity-api/service/subscriber_manager"
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

var SubscriberTypes map[string]subscriber_manager_pb.SubscriberType = map[string]subscriber_manager_pb.SubscriberType{
	"transmitter": subscriber_manager_pb.SubscriberType_TRANSMITTER,
	"exporter":    subscriber_manager_pb.SubscriberType_EXPORTER,
}

type MessageHandler func(*Message)

type Subscriber struct {
	client        *core.Client
	options       *Options
	host          string
	id            string
	pipelines     map[uint64]*Pipeline
	collectionMap map[string][]string
	subscription  *Subscription
	eventHandler  *EventHandler
	runner        *Runner
}

func NewSubscriber(options *Options) *Subscriber {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	subscriber := &Subscriber{
		options:       options,
		pipelines:     make(map[uint64]*Pipeline),
		collectionMap: make(map[string][]string),
		runner:        NewRunner(),
	}

	subscriber.subscription = NewSubscription(subscriber, options.BufferSize)
	subscriber.eventHandler = NewEventHandler(subscriber)

	return subscriber
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

	// Prepare token
	key := sub.options.Key
	token, err := key.Encryption().PrepareToken()
	if err != nil {
		return err
	}

	request := subscriber_manager_pb.RegisterSubscriberRequest{
		SubscriberID: subscriberID,
		Name:         name,
		Type:         subscriberType,
		Component:    component,
		AppID:        key.GetAppID(),
		Token:        token,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := sub.request("subscriber_manager.registerSubscriber", msg, true)
	if err != nil {
		return err
	}

	var reply subscriber_manager_pb.RegisterSubscriberReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	sub.id = subscriberID

	log.WithFields(logrus.Fields{
		"id": subscriberID,
	}).Info("Subscriber was registered")

	return nil
}

func (sub *Subscriber) healthCheck() error {

	// Fetch events from pipelines
	request := subscriber_manager_pb.HealthCheckRequest{
		SubscriberID: sub.id,
	}

	msg, _ := proto.Marshal(&request)

	respData, err := sub.request("subscriber_manager.healthCheck", msg, true)
	if err != nil {
		return err
	}

	var reply subscriber_manager_pb.HealthCheckReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		log.Error(reply.Reason)
		return err
	}

	return nil
}

func (sub *Subscriber) Connect(host string, options *core.Options) error {
	sub.client = core.NewClient()
	return sub.client.Connect(host, options)
}

func (sub *Subscriber) Disconnect() {
	sub.client.Disconnect()
}

func (sub *Subscriber) GetEndpoint() (*core.Endpoint, error) {
	return sub.client.ConnectToEndpoint(sub.options.Endpoint, sub.options.Domain, nil)
}

func (sub *Subscriber) Register(subscriberType subscriber_manager_pb.SubscriberType, component string, subscriberID string, name string) error {

	id := subscriberID
	if len(id) == 0 {
		id = uuid.NewV1().String()
	}

	// Getting endpoint from client object
	endpoint, err := sub.GetEndpoint()
	if err != nil {
		return err
	}

	// Register subscriber ID
	err = sub.register(subscriberType, component, id, name)
	if err != nil {
		return err
	}

	// Subscribe to channel
	channel := fmt.Sprintf("subscriber.%s", sub.id)

	log.WithFields(logrus.Fields{
		"channel": channel,
	}).Info("Subscribe to synchonizer")

	// Subscribe to channel
	conn := sub.client.GetConnection()
	s, err := conn.Subscribe(endpoint.Channel(channel), func(m *nats.Msg) {
		sub.eventHandler.ProcessEvent(m.Data)
		/*
			msg := messagePool.Get().(*Message)
			msg.Subscription = pipeline.subscriber.subscription
			msg.Pipeline = pipeline
			//		msg.Type = EVENT_TYPE
			msg.Data = m.Data
			msg.Callback = record

			sub.subscription.Push(m.Data)
		*/
		m.Ack()
	})
	if err != nil {
		return err
	}

	sub.subscription.sub = s
	sub.subscription.start()

	return nil
}

func (sub *Subscriber) Start() {

	// Initializing runner
	go sub.runner.Start()
	for _, pipeline := range sub.pipelines {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Info("Added pipeline to runner")
		sub.runner.AddPipeline(pipeline)
	}

	go func() {
		for {
			err := sub.healthCheck()
			if err != nil {
				log.Error(err)
			}

			<-time.After(time.Second * 30)
		}
	}()
}

func (sub *Subscriber) SetEventHandler(cb MessageHandler) {
	sub.subscription.eventHandler = cb
}

func (sub *Subscriber) SetSnapshotHandler(cb MessageHandler) {
	sub.subscription.snapshotHandler = cb
}

func (sub *Subscriber) GetPipelineCount() (uint64, error) {

	// Getting pipeline count
	opts := pipeline_manager.NewOptions()
	opts.Verbose = true
	opts.Key = sub.options.Key
	pm := pipeline_manager.NewPipelineManagerWithClient(sub.client, opts)
	return pm.GetPipelineCount()
}

func (sub *Subscriber) AddAllPipelines() error {
	return sub.SubscribeToPipelines(nil)
}

func (sub *Subscriber) SubscribeToPipelines(pipelines []uint64) error {

	// Getting pipeline count
	count, err := sub.GetPipelineCount()
	if err != nil {
		return err
	}

	if pipelines == nil {

		// Subscribe to all pipelines
		for i := uint64(0); i < count; i++ {
			pipeline := NewPipeline(sub, i, 0)
			err := sub.AddPipeline(pipeline)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// Subscribe to specific pipelines
	for _, pipelineID := range pipelines {
		if pipelineID >= count {
			return fmt.Errorf("No such pipeline: %d", pipelineID)
		}
	}

	for _, pipelineID := range pipelines {
		pipeline := NewPipeline(sub, pipelineID, 0)
		err := sub.AddPipeline(pipeline)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sub *Subscriber) AddPipeline(pipeline *Pipeline) error {

	if sub.options.StateStore != nil {
		pipelineState, err := sub.options.StateStore.GetPipelineState(pipeline.id)
		if err != nil {
			return err
		}

		// Load state
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
			"lastSeq":  pipelineState.GetLastSequence(),
		}).Info("Loaded pipeline state from store")
		pipeline.UpdateLastSequence(pipelineState.GetLastSequence())
	}

	// Initializing pipeline
	err := pipeline.Initialize()
	if err != nil {
		return err
	}

	sub.pipelines[pipeline.id] = pipeline

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
	opts := subscriber_manager.NewOptions()
	opts.Verbose = sub.options.Verbose
	opts.Endpoint = sub.options.Endpoint
	opts.Domain = sub.options.Domain
	opts.Key = sub.options.Key
	sm := subscriber_manager.NewSubscriberManagerWithClient(sub.client, opts)

	return sm.SubscribeToCollections(sub.id, collections)
}

func (sub *Subscriber) GetCollectionInfo(collection string) []string {
	tables, ok := sub.collectionMap[collection]
	if !ok {
		return nil
	}

	return tables
}

func (sub *Subscriber) GetPipeline(pipelineID uint64) *Pipeline {

	pipeline, ok := sub.pipelines[pipelineID]
	if ok {
		return pipeline
	}

	return nil
}

/*
func (sub *Subscriber) AwakePipeline(pipelineID uint64) {

	if sub.scheduler == nil {
		return
	}

	sub.scheduler.Awake(pipelineID)
}
*/
func (sub *Subscriber) ReleasePipeline(pipelineID uint64) {

	//	sub.scheduler.Idle(pipelineID)
}
