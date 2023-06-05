package biz

import (
	// "fmt"
	"time"

	"github.com/d2jvkpn/collector/proto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Record struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	Status            bool               `bson:"status" json:"status"`
	Tags              []string           `bson:"tags" json:"tags"`
	EventAt           time.Time          `bson:"eventAt" json:"eventAt"`
	*proto.RecordData `bson:",inline"`
}

/*
type Data struct {
	ServiceName    string            `bson:"serviceName" json:"serviceName"`
	ServiceVersion string            `bson:"serviceVersion,omitempty" json:"serviceVersion,omitempty"`
	EventId        string            `bson:"eventId" json:"eventId"`
	EventAt        time.Time         `bson:"eventAt" json:"eventAt"`
	BizName        string            `bson:"bizName" json:"bizName"`
	BizVersion     string            `bson:"bizVersion,omitempty" json:"bizVersion,omitempty"`
	BindId         string            `bson:"bindId,omitempty" json:"bindId,omitempty"` // deprecated, use BindIds instead
	BindIds        map[string]string `bson:"bindIds,omitempty" json:"bindIds,omitempty"`
	// Data           any       `bson:"data" json:"data"`
	Data json.RawMessage `bson:"data" json:"data"`
}
*/

type DataMsg struct {
	RecordData *proto.RecordData
	Fields     []zap.Field
}

func RecordFromData(rd *proto.RecordData, createdAt time.Time) Record {
	// ?? if createdAt.IsZero()
	return Record{
		CreatedAt:  createdAt,
		Status:     false,
		Tags:       []string{},
		EventAt:    rd.EventTimestamp.AsTime(),
		RecordData: rd,
	}
}
