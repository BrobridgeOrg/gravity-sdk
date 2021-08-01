package eventstore

import (
	"os"

	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type EventStore struct {
	client  *core.Client
	options *Options
}

func NewEventStore(options *Options) *EventStore {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	es := &EventStore{
		options: options,
	}

	return es
}

func NewEventStoreWithClient(client *core.Client, options *Options) *EventStore {

	es := NewEventStore(options)
	es.client = client

	return es
}

func (es *EventStore) Connect(host string, options *core.Options) error {
	es.client = core.NewClient()
	return es.client.Connect(host, options)
}

func (es *EventStore) Disconnect() {
	es.client.Disconnect()
}

func (es *EventStore) GetEndpoint() (*core.Endpoint, error) {
	return es.client.ConnectToEndpoint(es.options.Endpoint, es.options.Domain, nil)
}
