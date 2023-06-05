package internal

import (
	// "fmt"

	"github.com/d2jvkpn/collector/internal/biz"
	"github.com/d2jvkpn/collector/pkg/kafka"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	_Logger    *zap.Logger
	_Metrics   *viper.Viper
	_CloseOtel func() error

	_MongoClient   *mongo.Client
	_KafkaHandler  *kafka.Handler
	_RecordHandler *biz.RecordHandler
	_GrpcServer    *grpc.Server
)
