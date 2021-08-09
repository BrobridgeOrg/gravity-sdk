package gravity_sdk_types_record

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
	"unsafe"

	"github.com/golang/protobuf/ptypes"
)

var (
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
		return math.Float64frombits(binary.LittleEndian.Uint64(value.Value))
	case DataType_INT64:
		return int64(binary.LittleEndian.Uint64(value.Value))
	case DataType_UINT64:
		return uint64(binary.LittleEndian.Uint64(value.Value))
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

	default:
		data, _ := getBytes(data)
		return &Value{
			Type:  DataType_BINARY,
			Value: data,
		}, nil
	}
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
		binary.LittleEndian.PutUint64(buf, uint64(data.(uint)))
	case reflect.Uint8:
		binary.LittleEndian.PutUint64(buf, uint64(data.(uint8)))
	case reflect.Uint16:
		binary.LittleEndian.PutUint64(buf, uint64(data.(uint16)))
	case reflect.Uint32:
		binary.LittleEndian.PutUint64(buf, uint64(data.(uint32)))
	case reflect.Uint64:
		binary.LittleEndian.PutUint64(buf, data.(uint64))
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
		binary.LittleEndian.PutUint64(buf, uint64(data.(int)))
	case reflect.Int8:
		binary.LittleEndian.PutUint64(buf, uint64(data.(int8)))
	case reflect.Int16:
		binary.LittleEndian.PutUint64(buf, uint64(data.(int16)))
	case reflect.Int32:
		binary.LittleEndian.PutUint64(buf, uint64(data.(int32)))
	case reflect.Int64:
		binary.LittleEndian.PutUint64(buf, uint64(data.(int64)))
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
		binary.Write(&buf, binary.LittleEndian, data)
	case reflect.Float64:
		binary.Write(&buf, binary.LittleEndian, data)
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
