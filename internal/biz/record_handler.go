package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
	"github.com/valyala/fastjson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type RecordHandler struct {
	bp     *gotk.BatchProcess[*DataMsg]
	logger *zap.Logger
	db     *mongo.Database
}

func NewRecordHandler(vp *viper.Viper) (handler *RecordHandler, err error) {
	var (
		count    int
		interval time.Duration
	)

	count = vp.GetInt("count")
	if count <= 0 {
		return nil, fmt.Errorf("NewHandler: invalid count")
	}

	interval = vp.GetDuration("interval")
	if interval <= 0 {
		return nil, fmt.Errorf("NewHandler: invalid interval")
	}

	handler = &RecordHandler{}

	handler.bp, err = gotk.NewBatchProcess[*DataMsg](count, interval, handler.InsertMany)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (handler *RecordHandler) Down() {
	handler.bp.Down()
}

func (handler *RecordHandler) WithLogger(logger *zap.Logger) *RecordHandler {
	handler.logger = logger
	return handler
}

func (handler *RecordHandler) WithDatabase(db *mongo.Database) *RecordHandler {
	handler.db = db
	return handler
}

func (handler *RecordHandler) Ok() (err error) {
	if handler.bp == nil {
		return fmt.Errorf("bp is unset")
	}

	if handler.logger == nil {
		return fmt.Errorf("logger is unset")
	}

	if handler.db == nil {
		return fmt.Errorf("db is unset")
	}

	return nil
}

func (handler *RecordHandler) Handle(msg *sarama.ConsumerMessage) {
	var (
		err  error
		data DataMsg
	)

	data.Fields = fieldsFromMsg(msg)

	if err = json.Unmarshal(msg.Value, &data.RecordData); err != nil {
		handler.logger.Error(fmt.Sprintf("unmarshal RecordData: %v", err), data.Fields...)
		return
	}

	if err = fastjson.ValidateBytes(data.RecordData.Data); err != nil {
		handler.logger.Error(fmt.Sprintf("invalid RecordData.Data: %v", err), data.Fields...)
		return
	}

	if err = handler.bp.Recv(&data); err != nil {
		handler.logger.Error(fmt.Sprintf("unexpected recv: %v", err), data.Fields...)
		return
	}
}

func (handler *RecordHandler) InsertMany(dataList []*DataMsg) {
	var (
		err       error
		createdAt time.Time
		items     []any
		result    *mongo.InsertManyResult
	)

	if len(dataList) == 0 {
		return
	}

	createdAt = time.Now()
	items = make([]any, 0, len(dataList))

	for i := range dataList {
		items = append(items, RecordFromData(dataList[i].RecordData, createdAt))
	}

	at := createdAt.UTC()
	coll := fmt.Sprintf("records_%dS%d", at.Year(), (at.Month()+2)/3)

	result, err = handler.db.
		Collection(coll).
		InsertMany(context.TODO(), items)

	if err != nil {
		msgs := make([][]zap.Field, 0, len(items))
		for i := range dataList {
			msgs = append(msgs, dataList[i].Fields)
		}

		handler.logger.Error("insert_many", zap.Any("messages", msgs))
		return
	}

	handler.logger.Info(
		"insert_many",
		zap.Int("number", len(items)),
		zap.Any("result", result),
	)
}

func fieldsFromMsg(msg *sarama.ConsumerMessage) []zap.Field {
	return []zap.Field{
		zap.Time("timestamp", msg.Timestamp),
		zap.String("key", string(msg.Key)),
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
	}
}
