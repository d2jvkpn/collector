package proto

import (
	// "fmt"
	"encoding/json"
	// "time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewRecordData(serviceName, bizName string) *RecordData {
	return &RecordData{
		ServiceName:    serviceName,
		ServiceVersion: "0.1.0",
		// EventAt:        time.Now(),
		EventTimestamp: timestamppb.Now(),
		BizName:        bizName,
	}
}

// update ServiceVersion
func (data *RecordData) WithSvcV(version string) *RecordData {
	data.ServiceVersion = version
	return data
}

func (data *RecordData) WithEventId(eventId string) *RecordData {
	data.EventId = eventId
	return data
}

// update BizVersion
func (data *RecordData) WithBizV(version string) *RecordData {
	data.BizVersion = version
	return data
}

/*
// deprecated, use BindIds instead
func (data *Data) WithBindId(bindId string) *Data {
	data.BindId = bindId
	return data
}
*/

func (data *RecordData) WithBindIds(bindIds map[string]string) *RecordData {
	data.BindIds = bindIds
	return data
}

func (data *RecordData) WithData(item any) *RecordData {
	// TODO: handle error
	bts, _ := json.Marshal(item)
	data.Data = bts
	return data
}
