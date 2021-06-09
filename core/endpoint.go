package core

import (
	"fmt"
	"sync"

	nats "github.com/nats-io/nats.go"
)

type Endpoint struct {
	client  *Client
	name    string
	domain  string
	options *EndpointOptions
	isReady bool
	mutex   sync.Mutex
}

func NewEndpoint(client *Client, name string, domain string, options *EndpointOptions) *Endpoint {

	endpoint := &Endpoint{
		client:  client,
		name:    name,
		domain:  domain,
		options: options,
		isReady: false,
	}

	if options == nil {
		endpoint.options = NewEndpointOptions()
	}

	return endpoint
}

func (endpoint *Endpoint) Connect() error {

	endpoint.mutex.Lock()
	defer endpoint.mutex.Unlock()

	if !endpoint.isReady {
		//TODO: authentication

		endpoint.isReady = true
	}

	return nil
}

func (endpoint *Endpoint) GetConnection() *nats.Conn {
	return endpoint.client.eventbus.GetConnection()
}

func (endpoint *Endpoint) Channel(channelName string) string {
	return fmt.Sprintf("%s.%s", endpoint.domain, channelName)
}
