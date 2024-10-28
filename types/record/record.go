package record

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamp_pb "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrNotFoundKeyPath    = errors.New("Not found key path")
	ErrNotSupportedMethod = errors.New("Not supported method")
)

var (
	NotUnsignedIntegerErr = errors.New("Not unisgned integer")
	NotIntegerErr         = errors.New("Not integer")
	NotFloatErr           = errors.New("Not float")
)

const (
	PathTokenTypeKey   = 0
	PathTokenTypeIndex = 1
)

type PathToken struct {
	Type  int
	Value string
}

func ParsePath(input string) []PathToken {
	var result []PathToken
	start := 0

	for i := 0; i < len(input); i++ {
		char := input[i]
		if char == '.' || char == '[' {
			if start < i {
				result = append(result, PathToken{Type: PathTokenTypeKey, Value: input[start:i]})
			}
			start = i + 1
		} else if char == ']' {
			if start < i {
				result = append(result, PathToken{Type: PathTokenTypeIndex, Value: input[start:i]})
			}
			start = i + 1
		}
	}

	if start < len(input) {
		result = append(result, PathToken{Type: PathTokenTypeKey, Value: input[start:]})
	}

	return result
}

func GetValueData(value *Value) interface{} {
	return getValueData(value, false)
}

func NewRecord() *Record {
	return &Record{
		Payload: &Value{
			Type: DataType_MAP,
			Map:  &MapValue{},
		},
	}
}

func (record *Record) GetValueByPath(key string) (*Value, error) {

	v, err := GetValueByPath(record.Payload, key)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (record *Record) GetValueDataByPath(key string) (interface{}, error) {

	v, err := GetValueByPath(record.Payload, key)
	if err != nil {
		return nil, err
	}

	return getValueData(v, false), nil
}

func (record *Record) CalculateKey(fields []string) ([]byte, error) {

	keys := make([][]byte, len(fields))

	for i, keyPath := range fields {
		v, err := record.GetValueByPath(keyPath)
		if err != nil {
			return nil, ErrNotFoundKeyPath
		}

		data, _ := v.GetBytes()
		keys[i] = data
	}

	return bytes.Join(keys, []byte("_")), nil

}

func (record *Record) AsMap() map[string]interface{} {
	return ConvertFieldsToMap(record.Payload.Map.Fields)
}

func CreateValue(t DataType, data interface{}) (*Value, error) {

	if data == nil {
		return &Value{
			Type: DataType_NULL,
		}, nil
	}

	value := &Value{
		Type: t,
	}

	switch t {
	case DataType_INT64:
		data, _ := getBytesFromInteger(data)
		value.Value = data
	case DataType_UINT64:
		data, _ := getBytesFromUnsignedInteger(data)
		value.Value = data
	case DataType_FLOAT64:
		data, _ := getBytesFromFloat(data)
		value.Value = data
	case DataType_BOOLEAN:
		data, _ := getBytes(data)
		value.Value = data
	case DataType_STRING:
		value.Value = StrToBytes(data.(string))
	case DataType_TIME:
		value.Timestamp = timestamp_pb.New(data.(time.Time))
	case DataType_BINARY:
		value.Value = data.([]uint8)
	case DataType_MAP:

		// Prepare map value
		mv := &MapValue{
			Fields: make([]*Field, 0),
		}

		v := reflect.ValueOf(data)

		// Convert each key-value set
		for _, key := range v.MapKeys() {
			ele := v.MapIndex(key)

			// Convert value to protobuf format
			v, err := GetValueFromInterface(ele.Interface())
			if err != nil {
				fmt.Println(err)
				continue
			}

			field := Field{
				Name:  key.Interface().(string),
				Value: v,
			}

			mv.Fields = append(mv.Fields, &field)
		}

		value.Map = mv

	case DataType_ARRAY:

		v := reflect.ValueOf(data)

		// Prepare map value
		av := &ArrayValue{
			Elements: make([]*Value, 0, v.Len()),
		}

		for i := 0; i < v.Len(); i++ {
			ele := v.Index(i)

			// Convert value to protobuf format
			v, err := GetValueFromInterface(ele.Interface())
			if err != nil {
				fmt.Println(err)
				continue
			}

			av.Elements = append(av.Elements, v)
		}

		value.Array = av
	}

	return value, nil
}

func GetValueFromInterface(data interface{}) (*Value, error) {

	if data == nil {
		return &Value{
			Type: DataType_NULL,
		}, nil
	}

	switch data.(type) {
	case json.Number:

		if n, err := data.(json.Number).Int64(); err == nil {
			// Integer
			return CreateValue(DataType_INT64, n)
		} else if f, err := data.(json.Number).Float64(); err == nil {
			// Float
			return CreateValue(DataType_FLOAT64, f)
		}
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Int64:
		return CreateValue(DataType_INT64, data)
	case reflect.Uint64:
		return CreateValue(DataType_UINT64, data)
	case reflect.Float64:
		return CreateValue(DataType_FLOAT64, data)
	case reflect.Bool:
		return CreateValue(DataType_BOOLEAN, data)
	case reflect.String:
		return CreateValue(DataType_STRING, data)
	case reflect.Map:
		return CreateValue(DataType_MAP, data)
	case reflect.Slice:
		return CreateValue(DataType_ARRAY, data)
	}

	// Time
	switch d := data.(type) {
	case time.Time:
		return CreateValue(DataType_TIME, d)
	}

	// binary by default
	return CreateValue(DataType_BINARY, data)
}

func getBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func getBytesFromUnsignedInteger(data interface{}) ([]byte, error) {

	var buf = make([]byte, 8)

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Uint:
		binary.BigEndian.PutUint64(buf, uint64(data.(uint)))
	case reflect.Uint8:
		binary.BigEndian.PutUint64(buf, uint64(data.(uint8)))
	case reflect.Uint16:
		binary.BigEndian.PutUint64(buf, uint64(data.(uint16)))
	case reflect.Uint32:
		binary.BigEndian.PutUint64(buf, uint64(data.(uint32)))
	case reflect.Uint64:
		binary.BigEndian.PutUint64(buf, data.(uint64))
	default:
		return nil, NotUnsignedIntegerErr
	}

	return buf, nil
}

func getBytesFromInteger(data interface{}) ([]byte, error) {

	var buf = make([]byte, 8)

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Int:
		binary.BigEndian.PutUint64(buf, uint64(data.(int)))
	case reflect.Int8:
		binary.BigEndian.PutUint64(buf, uint64(data.(int8)))
	case reflect.Int16:
		binary.BigEndian.PutUint64(buf, uint64(data.(int16)))
	case reflect.Int32:
		binary.BigEndian.PutUint64(buf, uint64(data.(int32)))
	case reflect.Int64:
		binary.BigEndian.PutUint64(buf, uint64(data.(int64)))
	default:
		return nil, NotIntegerErr
	}

	return buf, nil
}

func getBytesFromFloat(data interface{}) ([]byte, error) {

	var buf bytes.Buffer

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Float32:
		binary.Write(&buf, binary.BigEndian, data)
	case reflect.Float64:
		binary.Write(&buf, binary.BigEndian, data)
	default:
		return nil, NotFloatErr
	}

	return buf.Bytes(), nil
}

func UnmarshalMapData(obj map[string]interface{}, record *Record) error {

	values := &MapValue{
		Fields: make([]*Field, 0, len(obj)),
	}

	for key, value := range obj {

		// Convert value to protobuf format
		v, err := GetValueFromInterface(value)
		if err != nil {
			fmt.Println(err)
			continue
		}

		values.Fields = append(values.Fields, &Field{
			Name:  key,
			Value: v,
		})
	}

	record.Payload = &Value{
		Type: DataType_MAP,
		Map:  values,
	}

	return nil
}

func UnmarshalJSON(data []byte, record *Record) error {

	var jsonObj map[string]interface{}
	err := jsoniter.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}

	return UnmarshalMapData(jsonObj, record)
}

func Unmarshal(data []byte, record *Record) error {
	return proto.Unmarshal(data, record)
}

func Marshal(record *Record) ([]byte, error) {
	return proto.Marshal(record)
}

func MarshalJSON(record *Record) ([]byte, error) {

	mapData := ConvertFieldsToMap(record.Payload.Map.Fields)

	return jsoniter.Marshal(mapData)
}

func GetField(fields []*Field, fieldName string) *Field {

	for _, field := range fields {
		if field.Name == fieldName {
			return field
		}
	}

	return nil
}

func GetValueByPath(value *Value, key string) (*Value, error) {

	var err error
	state := value

	tokens := ParsePath(key)
	for _, token := range tokens {
		if state, err = getValue(state, token); err != nil {
			return nil, err
		}
	}
	/*
		var err error
		state := value
		parser := parse(key)
		token := parser.nextItem()
		for token.typ != tokenEnd && token.typ != tokenError {
			if state, err = getValue(state, token); err != nil {
				return nil, err
			}
			token = parser.nextItem()
		}
		if token.typ == tokenError {
			return nil, fmt.Errorf(token.val)
		}
	*/
	return state, nil
}

func getValue(value *Value, token PathToken) (*Value, error) {
	switch value.Type {
	case DataType_MAP:
		switch token.Type {
		case PathTokenTypeKey:
			field := GetField(value.Map.Fields, token.Value)
			if field != nil {
				return field.Value, nil
			}
			return nil, fmt.Errorf("key not found")
		default:
			return nil, fmt.Errorf("key not found")
		}
	case DataType_ARRAY:
		switch token.Type {
		case PathTokenTypeIndex:
			index, err := strconv.Atoi(token.Value)
			if err != nil {
				return nil, fmt.Errorf("expected array index, but got %s", token.Value)
			}
			if index < 0 || index >= len(value.Array.Elements) {
				return nil, fmt.Errorf("index out of bounds %s", token.Value)
			}
			return value.Array.Elements[index], nil
		default:
			return nil, fmt.Errorf("key not found")
		}
	default:
		return nil, fmt.Errorf("can't deal with this type %d", value.Type)
	}
}

/*
func getValue(value *Value, key token) (*Value, error) {

		switch value.Type {
		case DataType_MAP:
			switch key.typ {
			case tokenIdentifier:
				field := GetField(value.Map.Fields, key.val)
				if field != nil {
					return field.Value, nil
				}

				return nil, fmt.Errorf("key not found")
			default:
				return nil, fmt.Errorf("key not found")
			}
		case DataType_ARRAY:
			switch key.typ {
			case tokenArrayIndex:
				index, err := strconv.Atoi(key.val)
				if err != nil {
					return nil, fmt.Errorf("expected array index, but got %s", key.val)
				}

				if index < 0 || index >= len(value.Array.Elements) {
					return nil, fmt.Errorf("index out of bounds %s", key.val)
				}
				return value.Array.Elements[index], nil
			default:
				return nil, fmt.Errorf("key not found")
			}
		default:
			return nil, fmt.Errorf("can't deal with this type %d", value.Type)

		}
	}
*/
func getValueData(value *Value, isFlat bool) interface{} {

	switch value.Type {
	case DataType_MAP:

		mapFields := value.Map.Fields
		payload := make(map[string]interface{}, len(mapFields))

		for _, mapField := range mapFields {
			payload[mapField.Name] = getValueData(mapField.Value, false)
		}

		if !isFlat {
			return payload
		}

		// Convert to JSON string
		data, _ := jsoniter.Marshal(payload)

		return data

	case DataType_ARRAY:

		arrayFields := value.Array.Elements
		payload := make([]interface{}, len(arrayFields))

		for i, arrayField := range arrayFields {
			payload[i] = getValueData(arrayField, false)
		}

		if !isFlat {
			return payload
		}

		// Convert to JSON string
		data, _ := jsoniter.Marshal(payload)

		return data

	case DataType_FLOAT64:
		return math.Float64frombits(binary.BigEndian.Uint64(value.Value))
	case DataType_INT64:
		return int64(binary.BigEndian.Uint64(value.Value))
	case DataType_UINT64:
		return uint64(binary.BigEndian.Uint64(value.Value))
	case DataType_BOOLEAN:
		return int8(value.Value[0]) & 1
	case DataType_STRING:
		return string(value.Value)
	case DataType_NULL:
		return nil
	case DataType_TIME:
		return value.Timestamp.AsTime()
	}

	// binary
	return value.Value
}

func ConvertFieldsToMap(fields []*Field) map[string]interface{} {

	payload := make(map[string]interface{}, len(fields))

	for _, field := range fields {
		value := getValueData(field.Value, false)
		payload[field.Name] = value
	}

	return payload
}

func ApplyChanges(orig *Value, changes *Value) {

	if orig == nil || changes == nil {
		return
	}

	if changes.Type != DataType_MAP {
		return
	}

	for _, field := range changes.Map.Fields {

		// Getting specifc field
		f := GetField(orig.Map.Fields, field.Name)
		if f == nil {

			// new field
			f = &Field{
				Name:  field.Name,
				Value: field.Value,
			}

			orig.Map.Fields = append(orig.Map.Fields, f)

			continue
		}

		// check type to update
		switch f.Value.Type {
		case DataType_ARRAY:
			f.Value.Array = field.Value.Array
		case DataType_MAP:
			ApplyChanges(f.Value, field.Value)
		default:
			// update value
			f.Value.Value = field.Value.Value
		}
	}

}

func Merge(origRecord *Record, updates *Record) []byte {

	if origRecord.Payload == nil {
		origRecord.Payload = &Value{
			Type: DataType_MAP,
			Map:  &MapValue{},
		}
	}

	// Update meta data
	origMeta := origRecord.Meta.AsMap()
	newMeta := updates.Meta.AsMap()

	for k, v := range newMeta {
		origMeta[k] = v
	}

	// Replace old meta data
	m, _ := structpb.NewStruct(origMeta)
	origRecord.Meta = m

	// Merge contents
	ApplyChanges(origRecord.Payload, updates.Payload)

	data, _ := Marshal(origRecord)

	return data
}
