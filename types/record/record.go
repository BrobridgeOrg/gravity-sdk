package gravity_sdk_types_record

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
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
