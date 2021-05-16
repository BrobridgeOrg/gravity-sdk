package gravity_sdk_types_projection

import (
	"bytes"
	"sync"

	gravity_sdk_types_record "github.com/BrobridgeOrg/gravity-sdk/types/record"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

type Field struct {
	Name    string      `json:"name"`
	Value   interface{} `json:"value"`
	Primary bool        `json:"primary"`
}

type Projection struct {
	EventName  string            `json:"event"`
	Collection string            `json:"collection"`
	Method     string            `json:"method"`
	PrimaryKey string            `json:"primaryKey"`
	Fields     []Field           `json:"fields"`
	Meta       map[string][]byte `json:"meta"`
}

type JSONResult struct {
	EventName  string                 `json:"event"`
	Collection string                 `json:"collection"`
	Payload    map[string]interface{} `json:"payload"`
	Meta       map[string][]byte      `json:"meta"`
}

var jsonLib = jsoniter.ConfigCompatibleWithStandardLibrary

var pool = sync.Pool{
	New: func() interface{} {
		return &JSONResult{}
	},
}

func Unmarshal(data []byte, pj *Projection) error {

	decoder := jsonLib.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err := decoder.Decode(pj)
	if err != nil {
		return nil
	}

	return nil
}

func (pj *Projection) ToJSON() ([]byte, error) {

	// Allocation
	result := pool.Get().(*JSONResult)
	result.EventName = pj.EventName
	result.Collection = pj.Collection
	result.Payload = make(map[string]interface{})
	result.Meta = pj.Meta

	for _, field := range pj.Fields {
		switch field.Value.(type) {
		case jsoniter.Number:

			if n, err := field.Value.(jsoniter.Number).Int64(); err == nil {
				result.Payload[field.Name] = n
			} else if f, err := field.Value.(jsoniter.Number).Float64(); err == nil {
				result.Payload[field.Name] = f
			}
		default:
			result.Payload[field.Name] = field.Value
		}
	}

	data, err := jsonLib.Marshal(result)

	pool.Put(result)

	return data, err
}

func (pj *Projection) ToRecord() (*gravity_sdk_types_record.Record, error) {

	// Convert projection to record
	var record gravity_sdk_types_record.Record
	err := ConvertToRecord(pj, &record)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &record, nil
}
