package adapter

import (
	"fmt"
	"os"

	"github.com/BrobridgeOrg/gravity-sdk/v2/core"

	jsoniter "github.com/json-iterator/go"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	domainEvent = "$GVT.%s.EVENT.%s"
)

var log = logrus.New()

type Message struct {
	EventName string `json:"event"`
	Payload   []byte `json:"payload"`
}

type PubAckFuture interface {
	Ok() <-chan *nats.PubAck
	Err() <-chan error
	Msg() *nats.Msg
}

type AdapterConnector struct {
	id      string
	client  *core.Client
	js      nats.JetStreamContext
	options *Options
}

func NewAdapterConnector(options *Options) *AdapterConnector {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	ac := &AdapterConnector{
		options: options,
	}

	return ac
}

func NewAdapterConnectorWithClient(client *core.Client, options *Options) *AdapterConnector {

	ac := NewAdapterConnector(options)
	ac.client = client
	js, _ := client.GetJetStream()
	ac.js = js

	return ac
}

func (ac *AdapterConnector) Connect(host string, options *core.Options) error {

	ac.client = core.NewClient()
	return ac.client.Connect(host, options)
}

func (ac *AdapterConnector) Disconnect() error {
	return ac.Disconnect()
}

func (ac *AdapterConnector) Publish(eventName string, payload []byte, meta map[string]string) (*nats.PubAck, error) {

	msg := &Message{
		EventName: eventName,
		Payload:   payload,
	}

	data, _ := json.Marshal(msg)

	// Prepare JetStream context
	js, err := ac.client.GetJetStream()
	if err != nil {
		return nil, err
	}

	// Publish
	subject := fmt.Sprintf(domainEvent, ac.options.Domain, eventName)

	m := &nats.Msg{
		Subject: subject,
		Data:    data,
	}

	if meta != nil {
		for k, v := range meta {
			m.Header.Add(k, v)
		}
	}

	return js.PublishMsg(m)

	/*
		ac.buffer.Push(&Request{
			EventName: eventName,
			Payload:   payload,
		})
	*/
}

func (ac *AdapterConnector) PublishAsync(eventName string, payload []byte, meta map[string]string) (nats.PubAckFuture, error) {

	msg := &Message{
		EventName: eventName,
		Payload:   payload,
	}

	data, _ := json.Marshal(msg)
	/*
		// Getting endpoint from client object
		js, err := ac.client.GetJetStream()
		if err != nil {
			return nil, err
		}
	*/
	// Publish
	subject := fmt.Sprintf(domainEvent, ac.options.Domain, eventName)

	m := &nats.Msg{
		Subject: subject,
		Data:    data,
	}

	if meta != nil {
		for k, v := range meta {
			m.Header.Add(k, v)
		}
	}

	return ac.js.PublishMsgAsync(m)
}

func (ac *AdapterConnector) PublishAsyncComplete() <-chan struct{} {
	return ac.js.PublishAsyncComplete()
}
