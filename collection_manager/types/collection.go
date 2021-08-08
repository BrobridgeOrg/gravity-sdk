package types

import (
	"encoding/json"
	"time"

	collection_manager_pb "github.com/BrobridgeOrg/gravity-api/service/collection_manager"
	"github.com/golang/protobuf/ptypes"
)

type Collection struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"dsc"`
	CreatedAt time.Time `json:"createdAt"`
}

func UnmarshalProto(p *collection_manager_pb.Collection) *Collection {

	createdAt, _ := ptypes.Timestamp(p.CreatedAt)

	collection := &Collection{
		ID:        p.CollectionID,
		Name:      p.Name,
		Desc:      p.Desc,
		CreatedAt: createdAt,
	}

	return collection
}

func MarshalProto(c *Collection) *collection_manager_pb.Collection {

	createdAt, _ := ptypes.TimestampProto(c.CreatedAt)

	return &collection_manager_pb.Collection{
		CollectionID: c.ID,
		Name:         c.Name,
		Desc:         c.Desc,
		CreatedAt:    createdAt,
	}
}

func Unmarshal(data []byte) (*Collection, error) {

	var collection Collection
	err := json.Unmarshal(data, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func NewCollection(id string) *Collection {
	return &Collection{
		ID: id,
	}
}

func (collection *Collection) ToBytes() []byte {
	data, _ := json.Marshal(collection)
	return data
}
