package product_event

import (
	"errors"

	record_type "github.com/BrobridgeOrg/gravity-sdk/v2/types/record"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNotFoundPrimaryKeyData = errors.New("Not found primary data")
)

func Marshal(pe *ProductEvent) ([]byte, error) {
	return proto.Marshal(pe)
}

func Unmarshal(data []byte, pe *ProductEvent) error {
	return proto.Unmarshal(data, pe)
}

func (pe *ProductEvent) ToBytes() ([]byte, error) {
	return proto.Marshal(pe)
}

func (pe *ProductEvent) GetContent() (*record_type.Record, error) {
	var r record_type.Record
	err := record_type.Unmarshal(pe.Data, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (pe *ProductEvent) SetContent(r *record_type.Record) error {
	data, err := record_type.Marshal(r)
	if err != nil {
		return err
	}

	pe.Data = data

	return nil
}
