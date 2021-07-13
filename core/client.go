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
		PingInterval:        time.Duration(options.PingInterval),
		MaxPingsOutstanding: options.MaxPingsOutstanding,
		MaxReconnects:       options.MaxReconnects,
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

func (client *Client) Disconnect() {
	client.eventbus.Close()
}

func (client *Client) ConnectToEndpoint(name string, domain string, options *EndpointOptions) (*Endpoint, error) {

	// Attempt to get existing endpoint
	endpoint := client.GetEndpoint(domain)
	if endpoint != nil {
		return endpoint, nil
	}

	// Create a new link to endpoint
	endpoint = NewEndpoint(client, name, domain, options)
	client.endpoints.Store(name, endpoint)

	err := endpoint.Connect()
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}

func (client *Client) GetEndpoint(name string) *Endpoint {

	v, ok := client.endpoints.Load(name)
	if ok {
		return v.(*Endpoint)
	}

	return nil
}

func (client *Client) GetConnection() *nats.Conn {
	return client.eventbus.GetConnection()
}
