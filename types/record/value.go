package record

import (
	"encoding/binary"
	"errors"

	jsoniter "github.com/json-iterator/go"
	//"google.golang.org/protobuf/proto"
)

var (
	UnsupportedOperationErr = errors.New("Unsupported operation")
)

func (v *Value) GetData() interface{} {
	return getValueData(v, false)
}

func (v *Value) GetBytes() ([]byte, error) {

	switch v.Type {
	case DataType_MAP:

		mapFields := v.Map.Fields
		payload := make(map[string]interface{}, len(mapFields))

		for _, mapField := range mapFields {
			payload[mapField.Name] = getValueData(mapField.Value, false)
		}

		// Convert to JSON string
		data, _ := jsoniter.Marshal(payload)
		//		data, _ := json.Marshal(payload)

		return data, nil

	case DataType_ARRAY:

		arrayFields := v.Array.Elements
		payload := make([]interface{}, len(arrayFields))

		for i, arrayField := range arrayFields {
			payload[i] = getValueData(arrayField, false)
		}

		// Convert to JSON string
		data, _ := jsoniter.Marshal(payload)
		//data, _ := json.Marshal(payload)

		return data, nil

	case DataType_TIME:
		ts := v.Timestamp.AsTime().UnixNano()
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(ts))
		return b, nil
	}

	return v.Value, nil
}
