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

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/proto"
)

var (
	NotFoundPrimaryKey    = errors.New("Not found primary key")
	NotUnsignedIntegerErr = errors.New("Not unisgned integer")
	NotIntegerErr         = errors.New("Not integer")
	NotFloatErr           = errors.New("Not float")
)

type RecordDef struct {
	HasPrimary    bool
	PrimaryColumn string
	Values        map[string]interface{}
	ColumnDefs    []*ColumnDef
}

type ColumnDef struct {
	ColumnName  string
	BindingName string
	Value       interface{}
}

func GetValue(value *Value) interface{} {

	switch value.Type {
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
		ts, _ := ptypes.Timestamp(value.Timestamp)
		return ts
	}

	// binary
	return value.Value
}

func GetDefinition(record *Record) (*RecordDef, error) {

	recordDef := &RecordDef{
		HasPrimary: false,
		Values:     make(map[string]interface{}),
		ColumnDefs: make([]*ColumnDef, 0, len(record.Fields)),
	}

	// Scanning fields
	for n, field := range record.Fields {

		value := GetValue(field.Value)

		// Primary key
		//		if field.IsPrimary == true {
		if record.PrimaryKey == field.Name {
			recordDef.Values["primary_val"] = value
			recordDef.HasPrimary = true
			recordDef.PrimaryColumn = field.Name
			continue
		}

		// Generate binding name
		bindingName := fmt.Sprintf("val_%s", strconv.Itoa(n))
		recordDef.Values[bindingName] = value

		// Store definition
		recordDef.ColumnDefs = append(recordDef.ColumnDefs, &ColumnDef{
			ColumnName:  field.Name,
			Value:       field.Name,
			BindingName: bindingName,
		})
	}

	if len(record.PrimaryKey) > 0 && !recordDef.HasPrimary {
		return nil, errors.New("Not found primary key")
	}

	return recordDef, nil
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
			bytes, err := getBytesFromInteger(n)
			if err == nil {
				return &Value{
					Type:  DataType_INT64,
					Value: bytes,
				}, nil
			}
		} else if f, err := data.(json.Number).Float64(); err == nil {
			// Float
			bytes, err := getBytesFromFloat(f)
			if err == nil {
				return &Value{
					Type:  DataType_FLOAT64,
					Value: bytes,
				}, nil
			}
		}
	}
	/*
		// Unsigned integer
		bytes, err := t.getBytesFromUnsignedInteger(data)
		if err == nil {
			return &transmitter.Value{
				Type:  transmitter.DataType_INT64,
				Value: bytes,
			}, nil
		}
	*/
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Int64:
		data, _ := getBytesFromInteger(data)
		return &Value{
			Type:  DataType_INT64,
			Value: data,
		}, nil
	case reflect.Uint64:
		data, _ := getBytesFromUnsignedInteger(data)
		return &Value{
			Type:  DataType_UINT64,
			Value: data,
		}, nil
	case reflect.Float64:
		data, _ := getBytesFromFloat(data)
		return &Value{
			Type:  DataType_FLOAT64,
			Value: data,
		}, nil
	case reflect.Bool:
		data, _ := getBytes(data)
		return &Value{
			Type:  DataType_BOOLEAN,
			Value: data,
		}, nil
	case reflect.String:
		return &Value{
			Type:  DataType_STRING,
			Value: StrToBytes(data.(string)),
		}, nil
	case reflect.Map:

		// Prepare map value
		value := MapValue{
			Fields: make([]*Field, 0),
		}

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

			value.Fields = append(value.Fields, &field)
		}

		return &Value{
			Type: DataType_MAP,
			Map:  &value,
		}, nil

	case reflect.Slice:

		// Prepare map value
		value := ArrayValue{
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

			value.Elements = append(value.Elements, v)
		}

		return &Value{
			Type:  DataType_ARRAY,
			Array: &value,
		}, nil
	}

	switch d := data.(type) {
	case time.Time:
		t, err := ptypes.TimestampProto(d)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return &Value{
			Type:      DataType_TIME,
			Timestamp: t,
		}, nil
	}

	// binary by default
	value, _ := getBytes(data)
	return &Value{
		Type:  DataType_BINARY,
		Value: value,
	}, nil
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

	record.Fields = make([]*Field, 0, len(obj))
	for key, value := range obj {

		// Convert value to protobuf format
		v, err := GetValueFromInterface(value)
		if err != nil {
			fmt.Println(err)
			continue
		}

		record.Fields = append(record.Fields, &Field{
			Name:  key,
			Value: v,
		})
	}

	return nil
}

func UnmarshalJSON(data []byte, record *Record) error {

	var jsonObj map[string]interface{}
	err := json.Unmarshal(data, &jsonObj)
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

func GetField(fields []*Field, fieldName string) *Field {

	for _, field := range fields {
		if field.Name == fieldName {
			return field
		}
	}

	return nil
}

func (record *Record) GetPayload() *Value {

	value := &Value{
		Type: DataType_MAP,
		Map: &MapValue{
			Fields: record.Fields,
		},
	}

	return value
}

func (record *Record) GetPrimaryKeyValue() (*Value, error) {

	value := &Value{
		Type: DataType_MAP,
		Map: &MapValue{
			Fields: record.Fields,
		},
	}

	v, err := GetValueByPath(value, record.PrimaryKey)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func GetValueByPath(value *Value, key string) (*Value, error) {
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
	return state, nil
}

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

func prepareValue(value *Value) interface{} {

	switch value.Type {
	case DataType_MAP:

		mapFields := value.Map.Fields
		payload := make(map[string]interface{}, len(mapFields))

		for _, mapField := range mapFields {

			payload[mapField.Name] = prepareValue(mapField.Value)
		}

		return payload

	case DataType_ARRAY:

		arrayFields := value.Array.Elements
		payload := make([]interface{}, len(arrayFields))

		for _, arrayField := range arrayFields {
			payload = append(payload, prepareValue(arrayField))
		}

		return payload

	default:
		return GetValue(value)
	}

	return GetValue(value)

}

func FieldsToMap(fields []*Field) map[string]interface{} {

	payload := make(map[string]interface{}, len(fields))

	for _, field := range fields {
		value := prepareValue(field.Value)
		payload[field.Name] = value
	}

	return payload
}
