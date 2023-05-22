package internal

import (
	// "fmt"

	"github.com/d2jvkpn/collector/pkg/kafka"
	"github.com/d2jvkpn/collector/pkg/wrap"

	"go.mongodb.org/mongo-driver/mongo"
	// "go.uber.org/zap"
)

var (
	_Logger       *wrap.Logger
	_MongoClient  *mongo.Client
	_KafkaHandler *kafka.Handler
	_Handler      *Handler
)
