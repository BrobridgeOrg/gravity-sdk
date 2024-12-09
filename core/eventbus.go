package core

import (
	"time"

	nats "github.com/nats-io/nats.go"
)

type EventBusOptions struct {
	PingInterval           time.Duration
	MaxPingsOutstanding    int
	MaxReconnects          int
	PublishAsyncMaxPending int
}

type EventBusHandler struct {
	Reconnect  func(natsConn *nats.Conn)
	Disconnect func(natsConn *nats.Conn)
}

type EventBus struct {
	connection *nats.Conn
	js         nats.JetStreamContext
	host       string
	handler    *EventBusHandler
	options    *EventBusOptions
}

func NewEventBus(host string, handler EventBusHandler, options EventBusOptions) *EventBus {
	return &EventBus{
		connection: nil,
		js:         nil,
		host:       host,
		handler:    &handler,
		options:    &options,
	}
}

func (eb *EventBus) Connect() error {

	nc, err := nats.Connect(eb.host,
		nats.RetryOnFailedConnect(true),
		nats.PingInterval(eb.options.PingInterval*time.Second),
		nats.MaxPingsOutstanding(eb.options.MaxPingsOutstanding),
		nats.MaxReconnects(eb.options.MaxReconnects),
		nats.ReconnectHandler(eb.ReconnectHandler),
		nats.DisconnectHandler(eb.handler.Disconnect),
	)
	if err != nil {
		return err
	}

	eb.connection = nc

	return nil
}
func (eb *EventBus) Close() {
	eb.connection.Close()
}
func (eb *EventBus) ReconnectHandler(natsConn *nats.Conn) {
	eb.handler.Reconnect(natsConn)
}

func (eb *EventBus) GetConnection() *nats.Conn {
	return eb.connection
}

func (eb *EventBus) GetJetStream() (nats.JetStreamContext, error) {

	if eb.js == nil {
		js, err := eb.connection.JetStream(nats.PublishAsyncMaxPending(eb.options.PublishAsyncMaxPending))
		if err != nil {
			return nil, err
		}

		eb.js = js
	}

	return eb.js, nil
}
