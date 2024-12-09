package core

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Client struct {
	host      string
	options   *Options
	eventbus  *EventBus
	endpoints sync.Map
}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) Connect(host string, options *Options) error {

	client.host = host
	client.options = options

	opts := EventBusOptions{
		PingInterval:           time.Duration(options.PingInterval),
		MaxPingsOutstanding:    options.MaxPingsOutstanding,
		MaxReconnects:          options.MaxReconnects,
		PublishAsyncMaxPending: options.PublishAsyncMaxPending,
	}

	// Create a new instance connector
	client.eventbus = NewEventBus(
		host,
		EventBusHandler{
			Reconnect: func(natsConn *nats.Conn) {
				// re-connected to event server
				client.options.ReconnectHandler()
			},
			Disconnect: func(natsConn *nats.Conn) {
				// event server was disconnected
				client.options.DisconnectHandler()
			},
		},
		opts,
	)

	// Connect to server
	err := client.eventbus.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) IsAuthenticationRequired() bool {
	if len(client.options.Token) > 0 {
		return true
	}

	return false
}

func (client *Client) Disconnect() {
	client.eventbus.Close()
}

func (client *Client) GetConnection() *nats.Conn {
	return client.eventbus.GetConnection()
}

func (client *Client) GetJetStream() (nats.JetStreamContext, error) {
	return client.eventbus.GetJetStream()
}

func (client *Client) Request(apiPath string, payload []byte, timeout time.Duration) (*nats.Msg, error) {

	msg := nats.NewMsg(apiPath)
	msg.Data = payload

	return client.RequestMsg(msg, timeout)
}

func (client *Client) RequestMsg(msg *nats.Msg, timeout time.Duration) (*nats.Msg, error) {

	if len(client.options.Token) > 0 {
		msg.Header.Add("Authorization", client.options.Token)
	}

	return client.GetConnection().RequestMsg(msg, timeout)
}
