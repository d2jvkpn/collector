package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/d2jvkpn/collector/models"

	"github.com/Shopify/sarama"
	"github.com/d2jvkpn/gotk/impls"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Handler struct {
	bp     *impls.BatchProcess[models.DataMsg]
	logger *zap.Logger
	db     *mongo.Database
}

func NewHandler(count int, duration time.Duration) (handler *Handler, err error) {
	handler = &Handler{}

	handler.bp, err = impls.NewBatchProcess[models.DataMsg](count, duration, handler.InsertMany)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (handler *Handler) Down() {
	handler.bp.Down()
}

func (handler *Handler) WithLogger(logger *zap.Logger) *Handler {
	handler.logger = logger
	return handler
}

func (handler *Handler) WithDatabase(db *mongo.Database) *Handler {
	handler.db = db
	return handler
}

func (handler *Handler) Ok() (err error) {
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

func (handler *Handler) Handle(msg *sarama.ConsumerMessage) {
	var (
		err  error
		data models.DataMsg
	)

	data.Fields = fieldsFromMsg(msg)

	if err = json.Unmarshal(msg.Value, &data.Data); err != nil {
		handler.logger.Error("failed to unmarshal msg", data.Fields...)
		return
	}

	if err = handler.bp.Recv(data); err != nil {
		handler.logger.Error(fmt.Sprintf("unexpected recv: %v", err), data.Fields...)
		return
	}
}

func (handler *Handler) InsertMany(dataList []models.DataMsg) {
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
		items = append(items, models.Record{
			CreatedAt: createdAt,
			Data:      dataList[i].Data,
		})
	}

	at := createdAt.UTC()
	coll := fmt.Sprintf("records_%dS%d", at.Year(), at.Month()%3)

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
