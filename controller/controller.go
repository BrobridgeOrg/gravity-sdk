package controller

import (
	"errors"
	"time"

	pb "github.com/BrobridgeOrg/gravity-api/service/controller"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/golang/protobuf/proto"
)

type Controller struct {
	client  *core.Client
	host    string
	options *Options
}

func NewController(options *Options) *Controller {
	return &Controller{
		options: options,
	}
}

func NewControllerWithClient(client *core.Client, options *Options) *Controller {
	return &Controller{
		client:  client,
		options: options,
	}
}

func (sub *Controller) Connect(host string, options *core.Options) error {
	sub.client = core.NewClient()
	return sub.client.Connect(host, options)
}

func (sub *Controller) Disconnect() {
	sub.client.Disconnect()
}

func (sub *Controller) SubscribeToCollections(subscriberID string, collections []string) error {

	conn := sub.client.GetConnection()

	request := pb.SubscribeToCollectionsRequest{
		SubscriberID: subscriberID,
		Collections:  collections,
	}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request("gravity.core.subscribeToCollections", msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply pb.SubscribeToCollectionsReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}
