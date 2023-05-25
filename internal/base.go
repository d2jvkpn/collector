package internal

import (
	// "fmt"

	"github.com/d2jvkpn/collector/pkg/kafka"
	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	_Config       *viper.Viper
	_Logger       *wrap.Logger
	_MongoClient  *mongo.Client
	_KafkaHandler *kafka.Handler
	_Handler      *Handler
)

func SetConfig(vp *viper.Viper) (err error) {
	_Config = vp

	// if _ServiceName = _Config.GetString("service_name"); _ServiceName == "" {
	// 	return fmt.Errorf("service_name is empty in config")
	// }

	return nil
}
