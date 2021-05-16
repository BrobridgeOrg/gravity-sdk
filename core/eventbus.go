package core

import (
	"time"

	nats "github.com/nats-io/nats.go"
)

type EventBusOptions struct {
	PingInterval        time.Duration
	MaxPingsOutstanding int
	MaxReconnects       int
}

type EventBusHandler struct {
	Reconnect  func(natsConn *nats.Conn)
	Disconnect func(natsConn *nats.Conn)
}

type EventBus struct {
	connection *nats.Conn
	host       string
	handler    *EventBusHandler
	options    *EventBusOptions
}

func NewEventBus(host string, handler EventBusHandler, options EventBusOptions) *EventBus {
	return &EventBus{
		connection: nil,
		host:       host,
		handler:    &handler,
		options:    &options,
	}
}

func (eb *EventBus) Connect() error {

	nc, err := nats.Connect(eb.host,
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
