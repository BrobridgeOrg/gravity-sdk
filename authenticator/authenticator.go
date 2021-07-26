package authenticator

import (
	"errors"
	"os"

	authenticator_pb "github.com/BrobridgeOrg/gravity-api/service/auth"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/BrobridgeOrg/gravity-sdk/core/encryption"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Authenticator struct {
	client     *core.Client
	options    *Options
	encryption *encryption.Encryption
}

func NewAuthenticator(options *Options) *Authenticator {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	// Initializing encryption
	en := encryption.NewEncryption()
	en.SetKey(options.AccessKey)

	return &Authenticator{
		options:    options,
		encryption: en,
	}
}

func NewAuthenticatorWithClient(client *core.Client, options *Options) *Authenticator {

	auth := NewAuthenticator(options)
	auth.client = client

	return auth
}

func (auth *Authenticator) Connect(host string, options *core.Options) error {
	auth.client = core.NewClient()
	return auth.client.Connect(host, options)
}

func (auth *Authenticator) Disconnect() {
	auth.client.Disconnect()
}

func (auth *Authenticator) GetEndpoint() (*core.Endpoint, error) {
	return auth.client.ConnectToEndpoint(auth.options.Endpoint, auth.options.Domain, nil)
}

func (auth *Authenticator) GetChannel() (string, error) {

	if len(auth.options.Channel) != 0 {
		return auth.options.Channel, nil
	}

	endpoint, err := auth.GetEndpoint()
	if err != nil {
		return "", err
	}

	return endpoint.Channel("auth"), nil
}

func (auth *Authenticator) Authenticate(appID string, token []byte) (*Entity, error) {

	// Prepare request
	request := authenticator_pb.AuthenticateRequest{
		AppID: appID,
		Token: token,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := auth.request("authenticate", msg, true)
	if err != nil {
		return nil, err
	}

	var reply authenticator_pb.AuthenticateReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	return ParseEntityProto(reply.Entity)
}

func (auth *Authenticator) CreateEntity(entity *Entity) error {

	// Prepare request
	request := authenticator_pb.CreateEntityRequest{
		Entity: ConvertEntityToProto(entity),
	}
	msg, _ := proto.Marshal(&request)

	respData, err := auth.request("createEntity", msg, true)
	if err != nil {
		return err
	}

	var reply authenticator_pb.CreateEntityReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (auth *Authenticator) UpdateEntity(entity *Entity) error {

	// Prepare request
	request := authenticator_pb.UpdateEntityRequest{
		AppID:  entity.AppID,
		Entity: ConvertEntityToProto(entity),
	}
	msg, _ := proto.Marshal(&request)

	respData, err := auth.request("createEntity", msg, true)
	if err != nil {
		return err
	}

	var reply authenticator_pb.UpdateEntityReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (auth *Authenticator) UpdateEntityKey(appID string, key string) error {

	// Prepare request
	request := authenticator_pb.UpdateEntityKeyRequest{
		AppID: appID,
		Key:   key,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := auth.request("updateEntityKey", msg, true)
	if err != nil {
		return err
	}

	var reply authenticator_pb.UpdateEntityKeyReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (auth *Authenticator) DeleteEntity(appID string) error {

	// Prepare request
	request := authenticator_pb.DeleteEntityRequest{
		AppID: appID,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := auth.request("deleteEntity", msg, true)
	if err != nil {
		return err
	}

	var reply authenticator_pb.DeleteEntityReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (auth *Authenticator) GetEntity(appID string) (*Entity, error) {

	// Prepare request
	request := authenticator_pb.GetEntityRequest{
		AppID: appID,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := auth.request("getEntity", msg, true)
	if err != nil {
		return nil, err
	}

	var reply authenticator_pb.GetEntityReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	return ParseEntityProto(reply.Entity)
}

func (auth *Authenticator) GetEntities(startID string, count int32) ([]*Entity, int32, error) {

	// Prepare request
	request := authenticator_pb.GetEntitiesRequest{
		StartID: startID,
		Count:   count,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := auth.request("getEntity", msg, true)
	if err != nil {
		return nil, 0, err
	}

	var reply authenticator_pb.GetEntitiesReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, 0, err
	}

	if !reply.Success {
		return nil, 0, errors.New(reply.Reason)
	}

	entities := make([]*Entity, 0)
	for _, entity := range reply.Entities {
		entry, _ := ParseEntityProto(entity)
		entities = append(entities, entry)
	}

	return entities, reply.Total, nil
}
