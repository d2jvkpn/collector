package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/d2jvkpn/collector/pkg/kafka"
	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Load(confg string) (err error) {
	var vp *viper.Viper

	vp = viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigFile(confg)
	if err = vp.ReadInConfig(); err != nil {
		return fmt.Errorf("ReadInConfig: %w", err)
	}
	vp.SetDefault("http.cors", "*")
	vp.SetDefault("log.size_mb", 256)

	_Logger, err = wrap.NewLogger(
		vp.GetString("log.path"),
		zap.InfoLevel,
		vp.GetInt("log.size_mb"),
	)
	if err != nil {
		return fmt.Errorf("NewLogger: %w", err)
	}

	if _MongoClient, err = wrap.MongoClient(vp, "mongodb"); err != nil {
		return fmt.Errorf("MongoClient: %w", err)
	}

	if _Handler, err = NewHandler(1000, time.Minute); err != nil {
		return fmt.Errorf("NewHandler: %w", err)
	}

	db := vp.GetString("mongodb.db")
	if db == "" {
		return fmt.Errorf("mongodb.db is empty")
	}

	database := _MongoClient.Database(db)
	err = _Handler.WithLogger(_Logger.Named("handler")).WithDatabase(database).Ok()
	if err != nil {
		return err
	}

	if _KafkaHandler, err = kafka.HandlerFromConfig(context.TODO(), vp, "kafka"); err != nil {
		return fmt.Errorf("HandlerFromConfig: %w", err)
	}
	_KafkaHandler.WithLogger(_Logger.Named("kafka")).WithHandle(_Handler.Handle)

	return nil
}
