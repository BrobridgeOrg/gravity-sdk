package record

import (
	"errors"

	"github.com/golang/protobuf/ptypes"
)

var (
	UnsupportedOperationErr = errors.New("Unsupported operation")
)

func (v *Value) GetBytes() ([]byte, error) {

	switch v.Type {
	case DataType_MAP:
		return nil, UnsupportedOperationErr
	case DataType_ARRAY:
		return nil, UnsupportedOperationErr
	case DataType_TIME:
		t, err := ptypes.Timestamp(v.Timestamp)
		if err != nil {
			return nil, err
		}

		return t.MarshalBinary()
	}

	return v.Value, nil
}
