package internal

import (
	// "fmt"

	"github.com/d2jvkpn/collector/pkg/kafka"

	"github.com/d2jvkpn/gotk/impls"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	_Config       *viper.Viper
	_Logger       *impls.Logger
	_MongoClient  *mongo.Client
	_KafkaHandler *kafka.Handler
	_Handler      *Handler
)
