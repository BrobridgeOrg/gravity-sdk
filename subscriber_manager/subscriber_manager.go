package subscriber_manager

import (
	"errors"
	"os"

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

	sm := &SubscriberManager{
		options: options,
	}

	return sm
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

	request := subscriber_manager_pb.GetSubscribersRequest{}
	msg, _ := proto.Marshal(&request)

	respData, err := sm.request("subscriber_manager.getSubscribers", msg, true)
	if err != nil {
		return nil, err
	}

	var reply subscriber_manager_pb.GetSubscribersReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	subscribers := make([]*Subscriber, len(reply.Subscribers))
	for i, sub := range reply.Subscribers {
		lastCheck, _ := ptypes.Timestamp(sub.LastCheck)
		subscribers[i] = &Subscriber{
			ID:          sub.SubscriberID,
			Name:        sub.Name,
			Component:   sub.Component,
			Type:        sub.Type,
			LastCheck:   lastCheck,
			AppID:       sub.AppID,
			AccessKey:   sub.AccessKey,
			Permissions: sub.Permissions,
			Collections: sub.Collections,
			Pipelines:   sub.Pipelines,
		}
	}

	return subscribers, nil
}

func (sm *SubscriberManager) SubscribeToCollections(subscriberID string, collections []string) error {

	request := subscriber_manager_pb.SubscribeToCollectionsRequest{
		SubscriberID: subscriberID,
		Collections:  collections,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := sm.request("subscriber_manager.subscribeToCollections", msg, true)
	if err != nil {
		return err
	}

	var reply subscriber_manager_pb.SubscribeToCollectionsReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}
