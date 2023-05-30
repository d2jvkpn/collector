package internal

import (
	// "fmt"

	"github.com/d2jvkpn/collector/internal/biz"
	"github.com/d2jvkpn/collector/pkg/kafka"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	_MongoClient  *mongo.Client
	_Handler      *biz.Handler
	_KafkaHandler *kafka.Handler
	_Logger       *zap.Logger
	_Metrics      *viper.Viper
)
