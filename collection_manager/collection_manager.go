package collection_manager

import (
	"errors"
	"os"

	collection_manager_pb "github.com/BrobridgeOrg/gravity-api/service/collection_manager"
	"github.com/BrobridgeOrg/gravity-sdk/collection_manager/types"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type CollectionManager struct {
	client  *core.Client
	options *Options
}

func NewCollectionManager(options *Options) *CollectionManager {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	return &CollectionManager{
		options: options,
	}
}

func NewCollectionManagerWithClient(client *core.Client, options *Options) *CollectionManager {

	cm := NewCollectionManager(options)
	cm.client = client

	return cm
}

func (cm *CollectionManager) Connect(host string, options *core.Options) error {
	cm.client = core.NewClient()
	return cm.client.Connect(host, options)
}

func (cm *CollectionManager) Disconnect() {
	cm.client.Disconnect()
}

func (cm *CollectionManager) GetEndpoint() (*core.Endpoint, error) {
	return cm.client.ConnectToEndpoint(cm.options.Endpoint, cm.options.Domain, nil)
}

func (cm *CollectionManager) Register(collection *types.Collection) error {

	request := collection_manager_pb.RegisterRequest{
		Collection: types.MarshalProto(collection),
	}
	msg, _ := proto.Marshal(&request)

	respData, err := cm.request("collection_manager.register", msg, true)
	if err != nil {
		return err
	}

	var reply collection_manager_pb.RegisterReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (cm *CollectionManager) Unregister(collectionID string) error {

	request := collection_manager_pb.UnregisterRequest{
		CollectionID: collectionID,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := cm.request("collection_manager.unregister", msg, true)
	if err != nil {
		return err
	}

	var reply collection_manager_pb.UnregisterReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (cm *CollectionManager) GetCollection(collectionID string) (*types.Collection, error) {

	request := collection_manager_pb.GetCollectionRequest{
		CollectionID: collectionID,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := cm.request("collection_manager.getCollections", msg, true)
	if err != nil {
		return nil, err
	}

	var reply collection_manager_pb.GetCollectionReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	return types.UnmarshalProto(reply.Collection), nil
}

func (cm *CollectionManager) GetCollections() ([]*types.Collection, error) {

	request := collection_manager_pb.GetCollectionsRequest{}
	msg, _ := proto.Marshal(&request)

	respData, err := cm.request("collection_manager.getCollections", msg, true)
	if err != nil {
		return nil, err
	}

	var reply collection_manager_pb.GetCollectionsReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	collections := make([]*types.Collection, len(reply.Collections))
	for i, c := range reply.Collections {
		collections[i] = types.UnmarshalProto(c)
	}

	return collections, nil
}
