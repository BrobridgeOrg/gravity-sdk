package adapter

import (
	"errors"
	"fmt"
	"os"
	"time"

	adapter_manager_pb "github.com/BrobridgeOrg/gravity-api/service/adapter_manager"
	dsa "github.com/BrobridgeOrg/gravity-api/service/dsa"
	"github.com/BrobridgeOrg/gravity-sdk/core"
	buffered_input "github.com/cfsghost/buffered-input"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type AdapterConnector struct {
	id      string
	client  *core.Client
	options *Options
	buffer  *buffered_input.BufferedInput
}

func NewAdapterConnector(options *Options) *AdapterConnector {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	ac := &AdapterConnector{
		options: options,
		//		buffer:  NewRequestBuffer(options.BatchSize),
	}

	// Initializing buffered input
	opts := buffered_input.NewOptions()
	opts.ChunkSize = 1000
	opts.ChunkCount = 1000
	opts.Timeout = 50 * time.Millisecond
	opts.Handler = ac.chunkHandler
	ac.buffer = buffered_input.NewBufferedInput(opts)

	return ac
}

func NewAdapterConnectorWithClient(client *core.Client, options *Options) *AdapterConnector {

	subscriber := NewAdapterConnector(options)
	subscriber.client = client

	return subscriber
}

func (ac *AdapterConnector) register(component string, adapterID string, name string) error {

	log.WithFields(logrus.Fields{
		"id": adapterID,
	}).Info("Registering adapter")

	conn := ac.client.GetConnection()

	// Attempt to connect to endpoint
	endpoint, err := ac.client.ConnectToEndpoint(ac.options.Endpoint, ac.options.Domain, nil)
	if err != nil {
		return err
	}

	request := adapter_manager_pb.RegisterAdapterRequest{
		AdapterID: adapterID,
		Name:      name,
		Component: component,
	}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request(endpoint.Channel("adapter_manager.register"), msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply adapter_manager_pb.RegisterAdapterReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	ac.id = adapterID

	log.WithFields(logrus.Fields{
		"id": adapterID,
	}).Info("Adapter was registered")

	return nil
}

func (ac *AdapterConnector) publish(requests []*Request) error {

	var done int32 = 0
	for {
		success, count, err := ac.BatchPublish(requests)
		if success {
			return nil
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		for i := int32(0); i < count; i++ {
			requests[done+i].IsCompleted = true
			done++
		}
	}
}

func (ac *AdapterConnector) chunkHandler(chunk []interface{}) {

	requests := make([]*Request, 0, len(chunk))
	for _, request := range chunk {
		requests = append(requests, request.(*Request))
	}

	ac.publish(requests)
}

func (ac *AdapterConnector) Connect(host string, options *core.Options) error {

	ac.client = core.NewClient()
	return ac.client.Connect(host, options)
}

func (ac *AdapterConnector) Disconnect() error {
	return ac.Disconnect()
}

func (ac *AdapterConnector) Release() {
	ac.buffer.Close()
}

func (ac *AdapterConnector) GetEndpoint() (*core.Endpoint, error) {
	return ac.client.ConnectToEndpoint(ac.options.Endpoint, ac.options.Domain, nil)
}

func (ac *AdapterConnector) Register(component string, adapterID string, name string) error {

	id := adapterID
	if len(id) == 0 {
		id = uuid.NewV1().String()
	}

	err := ac.register(component, id, name)
	if err != nil {
		return err
	}

	return nil
}

func (ac *AdapterConnector) BatchPublish(requests []*Request) (bool, int32, error) {

	request := &dsa.BatchPublishRequest{
		Requests: make([]*dsa.PublishRequest, 0, len(requests)),
	}

	for _, req := range requests {
		if req.IsCompleted {
			continue
		}

		request.Requests = append(request.Requests, &dsa.PublishRequest{
			EventName: req.EventName,
			Payload:   req.Payload,
		})
	}

	// Getting endpoint from client object
	endpoint, err := ac.GetEndpoint()
	if err != nil {
		return false, 0, err
	}

	connection := endpoint.GetConnection()

	// Sned
	reqMsg, _ := proto.Marshal(request)
	resp, err := connection.Request(endpoint.Channel("dsa.batch"), reqMsg, time.Second*30)
	if err != nil {
		return false, 0, err
	}

	// Parse reply message
	var reply dsa.BatchPublishReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return false, 0, err
	}

	if reply.Success {
		return true, reply.SuccessCount, nil
	}

	return false, reply.SuccessCount, errors.New(reply.Reason)
}

func (ac *AdapterConnector) Publish(eventName string, payload []byte, meta map[string]interface{}) error {

	ac.buffer.Push(&Request{
		EventName: eventName,
		Payload:   payload,
	})
	/*

		request := &dsa.PublishRequest{
			EventName: eventName,
			Payload:   payload,
		}

		reqMsg, _ := proto.Marshal(request)

		// Send
		connection := ac.eventbus.GetConnection()
		resp, err := connection.Request("gravity.dsa.incoming", reqMsg, time.Second*5)
		if err != nil {
			return err
		}

		// Parse reply message
		var reply dsa.PublishReply
		err = proto.Unmarshal(resp.Data, &reply)
		if err != nil {
			return err
		}

		if !reply.Success {
			return errors.New(reply.Reason)
		}
	*/
	return nil
}
