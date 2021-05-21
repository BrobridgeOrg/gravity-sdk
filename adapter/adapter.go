package adapter

import (
	"errors"
	"fmt"
	"os"
	"time"

	dsa "github.com/BrobridgeOrg/gravity-api/service/dsa"
	"github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type AdapterConnector struct {
	client  *core.Client
	options *Options
	buffer  *RequestBuffer
}

func NewAdapterConnector(options *Options) *AdapterConnector {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	ac := &AdapterConnector{
		options: options,
		buffer:  NewRequestBuffer(options.BatchSize),
	}

	go ac.startPublisher()

	return ac
}

func NewAdapterConnectorWithClient(client *core.Client, options *Options) *AdapterConnector {

	subscriber := NewAdapterConnector(options)
	subscriber.client = client

	return subscriber
}

func (ac *AdapterConnector) startPublisher() {

	for {
		select {
		case requests := <-ac.buffer.output:
			ac.publish(requests)
		case <-time.After(50 * time.Millisecond):
			ac.buffer.Flush()
		}
	}
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

func (ac *AdapterConnector) Connect(host string, options *core.Options) error {

	ac.client = core.NewClient()
	return ac.client.Connect(host, options)
}

func (ac *AdapterConnector) Disconnect() error {
	return ac.Disconnect()
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

	// Send
	connection := ac.client.GetConnection()
	reqMsg, _ := proto.Marshal(request)
	resp, err := connection.Request("gravity.dsa.batch", reqMsg, time.Second*30)
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
