package adapter

import (
	"errors"
	"time"

	dsa "github.com/BrobridgeOrg/gravity-api/service/dsa"
	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
)

type AdapterConnector struct {
	host     string
	options  *Options
	eventbus *EventBus
}

func NewAdapterConnector() *AdapterConnector {
	return &AdapterConnector{}
}

func (dc *AdapterConnector) Connect(host string, options *Options) error {

	dc.host = host
	dc.options = options

	opts := EventBusOptions{
		PingInterval:        time.Duration(options.PingInterval),
		MaxPingsOutstanding: options.MaxPingsOutstanding,
		MaxReconnects:       options.MaxReconnects,
	}

	// Create a new instance connector
	dc.eventbus = NewEventBus(
		host,
		EventBusHandler{
			Reconnect: func(natsConn *nats.Conn) {
				// re-connected to event server
				dc.options.ReconnectHandler()
			},
			Disconnect: func(natsConn *nats.Conn) {
				// event server was disconnected
				dc.options.DisconnectHandler()
			},
		},
		opts,
	)

	// Connect to server
	err := dc.eventbus.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (dc *AdapterConnector) Disconnect() error {
	return dc.Disconnect()
}

func (dc *AdapterConnector) Publish(eventName string, payload []byte, meta map[string]interface{}) error {

	request := &dsa.PublishRequest{
		EventName: eventName,
		Payload:   payload,
	}

	reqMsg, _ := proto.Marshal(request)

	// Send
	connection := dc.eventbus.GetConnection()
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

	return nil
}
