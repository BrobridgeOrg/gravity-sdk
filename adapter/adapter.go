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

// NewAdapterConnector creates a new instance of AdapterConnector using the provided options.
// This function initializes the AdapterConnector with the necessary configurations defined in Options.
// options: Configuration options for creating the AdapterConnector instance.
// Returns a pointer to the newly created AdapterConnector.
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

// NewAdapterConnectorWithClient creates a new instance of AdapterConnector with an existing core.Client and provided options.
// This function allows for more control and customization by utilizing an existing client instance.
// client: An existing core.Client instance to be used by the AdapterConnector.
// options: Configuration options for creating the AdapterConnector instance.
// Returns a pointer to the newly created AdapterConnector.
func NewAdapterConnectorWithClient(client *core.Client, options *Options) *AdapterConnector {

	ac := NewAdapterConnector(options)
	ac.client = client
	js, _ := client.GetJetStream()
	ac.js = js

	return ac
}

// Connect establishes a connection to the specified Gravity host using the provided options.
// host: The address of the Gravity service.
// options: Configuration options for the connection.
// Returns an error if the connection attempt fails.
func (ac *AdapterConnector) Connect(host string, options *core.Options) error {

	ac.client = core.NewClient()
	return ac.client.Connect(host, options)
}

// Disconnect closes the connection established by the AdapterConnector with the Gravity service.
// This function is used to cleanly terminate the connection before shutting down the application or when the connection is no longer needed.
// Returns an error if the disconnection process encounters any issues.
func (ac *AdapterConnector) Disconnect() error {
	return ac.Disconnect()
}

// Publish sends a message to the specified event name on Gravity.
// eventName: The name of the event to which the message will be published.
// payload: The byte array containing the message payload.
// meta: A map of metadata key-value pairs to be sent along with the message.
// Returns a PubAck which acknowledges the publishing of the message, and an error if the publishing fails.
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
		Header:  nats.Header{},
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

// PublishAsync sends a message asynchronously to the specified event name on Gravity.
// eventName: The name of the event to which the message will be published.
// payload: The byte array containing the message payload.
// meta: A map of metadata key-value pairs to be sent along with the message.
// Returns a PubAckFuture which can be used to receive the acknowledgment of the message once it's published, and an error if the publishing fails.
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
		m.Header = make(map[string][]string)
		for k, v := range meta {
			m.Header.Add(k, v)
		}
	}

	return ac.js.PublishMsgAsync(m)
}

// PublishAsyncComplete returns a channel that will be closed when all asynchronous publish operations are completed.
// This can be used to ensure all messages have been published before proceeding.
func (ac *AdapterConnector) PublishAsyncComplete() <-chan struct{} {
	return ac.js.PublishAsyncComplete()
}
