package internal

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

func (data *Data) WithSvcV(v string) *Data {
	data.ServiceVersion = v
	return data
}

func (data *Data) WithEventId(eid string) *Data {
	data.EventId = eid
	return data
}

func (data *Data) WithBizV(v string) *Data {
	data.BizVersion = v
	return data
}

func (data *Data) WithBindId(bId string) *Data {
	data.BindId = bId
	return data
}

func (data *Data) WithData(d any) *Data {
	data.Data = d
	return data
}
