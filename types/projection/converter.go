package gravity_sdk_types_projection

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"reflect"

	gravity_sdk_types_record "github.com/BrobridgeOrg/gravity-sdk/types/record"
	log "github.com/sirupsen/logrus"
)

var (
	NotUnsignedIntegerErr = errors.New("Not unisgned integer")
	NotIntegerErr         = errors.New("Not integer")
	NotFloatErr           = errors.New("Not float")
)

func ConvertToRecord(pj *Projection, record *gravity_sdk_types_record.Record) error {

	record.EventName = pj.EventName
	record.PrimaryKey = pj.PrimaryKey
	record.Fields = make([]*gravity_sdk_types_record.Field, 0, len(pj.Fields))
	/*
		record := &transmitter.Record{
			EventName: pj.EventName,
			Table:     table,
			Fields:    make([]*transmitter.Field, 0, len(pj.Fields)),
		}
	*/
	if pj.Method == "delete" {
		record.Method = gravity_sdk_types_record.Method_DELETE
	} else if pj.Method == "update" {
		record.Method = gravity_sdk_types_record.Method_UPDATE
	} else if pj.Method == "truncate" {
		record.Method = gravity_sdk_types_record.Method_TRUNCATE
	} else {
		record.Method = gravity_sdk_types_record.Method_INSERT
	}

	for _, field := range pj.Fields {

		// Convert value to protobuf format
		v, err := getValue(field.Value)
		if err != nil {
			log.Error(err)
			continue
		}

		record.Fields = append(record.Fields, &gravity_sdk_types_record.Field{
			Name:  field.Name,
			Value: v,
			//			IsPrimary: field.Primary,
		})
	}

	return nil
}

func getValue(data interface{}) (*gravity_sdk_types_record.Value, error) {

	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	switch data.(type) {
	case json.Number:

		if n, err := data.(json.Number).Int64(); err == nil {
			// Integer
			bytes, err := getBytesFromInteger(n)
			if err == nil {
				return &gravity_sdk_types_record.Value{
					Type:  gravity_sdk_types_record.DataType_INT64,
					Value: bytes,
				}, nil
			}
		} else if f, err := data.(json.Number).Float64(); err == nil {
			// Float
			bytes, err := getBytesFromFloat(f)
			if err == nil {
				return &gravity_sdk_types_record.Value{
					Type:  gravity_sdk_types_record.DataType_FLOAT64,
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
	case reflect.Bool:
		data, _ := getBytes(data)
		return &gravity_sdk_types_record.Value{
			Type:  gravity_sdk_types_record.DataType_BOOLEAN,
			Value: data,
		}, nil
	case reflect.String:
		return &gravity_sdk_types_record.Value{
			Type:  gravity_sdk_types_record.DataType_STRING,
			Value: StrToBytes(data.(string)),
		}, nil
	case reflect.Map:

		// Prepare map value
		value := gravity_sdk_types_record.MapValue{
			Fields: make([]*gravity_sdk_types_record.Field, 0),
		}

		// Convert each key-value set
		for _, key := range v.MapKeys() {
			ele := v.MapIndex(key)

			// Convert value to protobuf format
			v, err := getValue(ele.Interface())
			if err != nil {
				log.Error(err)
				continue
			}

			field := gravity_sdk_types_record.Field{
				Name:  key.Interface().(string),
				Value: v,
			}

			value.Fields = append(value.Fields, &field)
		}

		return &gravity_sdk_types_record.Value{
			Type: gravity_sdk_types_record.DataType_MAP,
			Map:  &value,
		}, nil

	case reflect.Slice:

		// Prepare map value
		value := gravity_sdk_types_record.ArrayValue{
			Elements: make([]*gravity_sdk_types_record.Value, 0, v.Len()),
		}

		for i := 0; i < v.Len(); i++ {
			ele := v.Index(i)

			// Convert value to protobuf format
			v, err := getValue(ele.Interface())
			if err != nil {
				log.Error(err)
				continue
			}

			value.Elements = append(value.Elements, v)
		}

		return &gravity_sdk_types_record.Value{
			Type:  gravity_sdk_types_record.DataType_ARRAY,
			Array: &value,
		}, nil

	default:
		data, _ := getBytes(data)
		return &gravity_sdk_types_record.Value{
			Type:  gravity_sdk_types_record.DataType_BINARY,
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
