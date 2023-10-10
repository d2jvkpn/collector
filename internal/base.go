package internal

import (
	// "fmt"

	"github.com/d2jvkpn/collector/internal/biz"

	"github.com/d2jvkpn/gotk/kafka"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	_Logger    *zap.Logger
	_Metrics   *viper.Viper
	_CloseOtel func() error

	_MongoClient   *mongo.Client
	_KafkaHandler  *kafka.Handler
	_RecordHandler *biz.RecordHandler
	_GrpcSS        *biz.GrpcServiceServer
)
