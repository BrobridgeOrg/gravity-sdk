package subscriber_manager

import (
	"errors"
	"os"
	"time"

	subscriber_manager_pb "github.com/BrobridgeOrg/gravity-api/service/subscriber_manager"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type SubscriberManager struct {
	client  *core.Client
	options *Options
}

func NewSubscriberManager(options *Options) *SubscriberManager {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	return &SubscriberManager{
		options: options,
	}
}

func NewSubscriberManagerWithClient(client *core.Client, options *Options) *SubscriberManager {

	subscriber := NewSubscriberManager(options)
	subscriber.client = client

	return subscriber
}

func (sm *SubscriberManager) Connect(host string, options *core.Options) error {
	sm.client = core.NewClient()
	return sm.client.Connect(host, options)
}

func (sm *SubscriberManager) Disconnect() {
	sm.client.Disconnect()
}

func (sm *SubscriberManager) GetEndpoint() (*core.Endpoint, error) {
	return sm.client.ConnectToEndpoint(sm.options.Endpoint, sm.options.Domain, nil)
}

func (sm *SubscriberManager) GetSubscribers() ([]*Subscriber, error) {

	// Getting endpoint from client object
	endpoint, err := sm.GetEndpoint()
	if err != nil {
		return nil, err
	}

	conn := endpoint.GetConnection()

	request := subscriber_manager_pb.GetSubscribersRequest{}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request(endpoint.Channel("subscriber_manager.getSubscribers"), msg, time.Second*10)
	if err != nil {
		return nil, err
	}

	var reply subscriber_manager_pb.GetSubscribersReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	subscribers := make([]*Subscriber, 0, len(reply.Subscribers))
	for _, sub := range reply.Subscribers {
		lastCheck, _ := ptypes.Timestamp(sub.LastCheck)
		subscribers = append(subscribers, &Subscriber{
			ID:        sub.SubscriberID,
			Name:      sub.Name,
			Component: sub.Component,
			Type:      sub.Type,
			LastCheck: lastCheck,
			Token:     sub.Token,
		})
	}

	return subscribers, nil
}

func (sm *SubscriberManager) SubscribeToCollections(subscriberID string, collections []string) error {

	// Getting endpoint from client object
	endpoint, err := sm.GetEndpoint()
	if err != nil {
		return err
	}

	conn := endpoint.GetConnection()

	request := subscriber_manager_pb.SubscribeToCollectionsRequest{
		SubscriberID: subscriberID,
		Collections:  collections,
	}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request(endpoint.Channel("subscriber_manager.subscribeToCollections"), msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply subscriber_manager_pb.SubscribeToCollectionsReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}
