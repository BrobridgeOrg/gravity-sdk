package gravity_sdk_types_projection

import (
	gravity_sdk_types_record "github.com/BrobridgeOrg/gravity-sdk/types/record"
	log "github.com/sirupsen/logrus"
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
		v, err := gravity_sdk_types_record.GetValueFromInterface(field.Value)
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
