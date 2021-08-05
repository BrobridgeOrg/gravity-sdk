package adapter_manager

import (
	"errors"
	"os"

	adapter_manager_pb "github.com/BrobridgeOrg/gravity-api/service/adapter_manager"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type AdapterManager struct {
	client  *core.Client
	options *Options
}

func NewAdapterManager(options *Options) *AdapterManager {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	return &AdapterManager{
		options: options,
	}
}

func NewAdapterManagerWithClient(client *core.Client, options *Options) *AdapterManager {

	adapter := NewAdapterManager(options)
	adapter.client = client

	return adapter
}

func (am *AdapterManager) Connect(host string, options *core.Options) error {
	am.client = core.NewClient()
	return am.client.Connect(host, options)
}

func (am *AdapterManager) Disconnect() {
	am.client.Disconnect()
}

func (am *AdapterManager) GetEndpoint() (*core.Endpoint, error) {
	return am.client.ConnectToEndpoint(am.options.Endpoint, am.options.Domain, nil)
}

func (am *AdapterManager) Unregister(adapterID string) error {

	request := adapter_manager_pb.UnregisterAdapterRequest{
		AdapterID: adapterID,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := am.request("adapter_manager.unregister", msg, true)
	if err != nil {
		return err
	}

	var reply adapter_manager_pb.UnregisterAdapterReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (am *AdapterManager) GetAdapters() ([]*Adapter, error) {

	request := adapter_manager_pb.GetAdaptersRequest{}
	msg, _ := proto.Marshal(&request)

	respData, err := am.request("adapter_manager.getAdapters", msg, true)
	if err != nil {
		return nil, err
	}

	var reply adapter_manager_pb.GetAdaptersReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	adapters := make([]*Adapter, len(reply.Adapters))
	for i, sub := range reply.Adapters {
		adapters[i] = &Adapter{
			ID:        sub.AdapterID,
			Name:      sub.Name,
			Component: sub.Component,
		}
	}

	return adapters, nil
}
