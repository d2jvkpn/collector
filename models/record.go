package models

import (
	// "fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Record struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Status    bool               `bson:"status" json:"status"`
	Tags      []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Data      `bson:",inline"`
}

type Data struct {
	ServiceName    string    `bson:"serviceName" json:"serviceName"`
	ServiceVersion string    `bson:"serviceVersion,omitempty" json:"serviceVersion,omitempty"`
	EventId        string    `bson:"eventId" json:"eventId"`
	EventAt        time.Time `bson:"eventAt" json:"eventAt"`
	BizName        string    `bson:"bizName" json:"bizName"`
	BizVersion     string    `bson:"bizVersion,omitempty" json:"bizVersion,omitempty"`
	BindId         string    `bson:"bindId,omitempty" json:"bindId,omitempty"`
	Data           any       `bson:"data" json:"data"`
}

type DataMsg struct {
	Data
	Fields []zap.Field
}

func NewData(serviceName, bizName string) *Data {
	return &Data{
		ServiceName:    serviceName,
		ServiceVersion: "0.1.0",
		EventAt:        time.Now(),
		BizName:        bizName,
	}
}

// update ServiceVersion
func (data *Data) WithSvcV(version string) *Data {
	data.ServiceVersion = version
	return data
}

func (data *Data) WithEventId(eventId string) *Data {
	data.EventId = eventId
	return data
}

// update BizVersion
func (data *Data) WithBizV(version string) *Data {
	data.BizVersion = version
	return data
}

func (data *Data) WithBindId(bindId string) *Data {
	data.BindId = bindId
	return data
}

func (item *Data) WithData(data any) *Data {
	item.Data = data
	return item
}
