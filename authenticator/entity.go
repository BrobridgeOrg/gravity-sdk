package auth

import (
	"encoding/json"

	authenticator_pb "github.com/BrobridgeOrg/gravity-api/service/auth"
)

type Entity struct {
	AppID      string
	AppName    string
	AccessKey  string
	Properties map[string]interface{}
}

func NewEntity() *Entity {
	return &Entity{
		Properties: make(map[string]interface{}),
	}
}

func ParseEntityProto(pbEntity *authenticator_pb.Entity) (*Entity, error) {

	entity := NewEntity()

	entity.AppID = pbEntity.AppID
	entity.AppName = pbEntity.AppName
	entity.AccessKey = pbEntity.Key

	err := json.Unmarshal(pbEntity.Properties, &entity.Properties)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func ConvertEntityToProto(entity *Entity) *authenticator_pb.Entity {

	props, _ := json.Marshal(entity.Properties)

	return &authenticator_pb.Entity{
		AppID:      entity.AppID,
		AppName:    entity.AppName,
		Key:        entity.AccessKey,
		Properties: props,
	}
}
