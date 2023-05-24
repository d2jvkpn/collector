package internal

import (
	"context"
	"fmt"
	// "time"

	"github.com/d2jvkpn/collector/pkg/kafka"
	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/d2jvkpn/gotk/impls"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Load(config string) (err error) {
	var vp *viper.Viper

	defer func() {
		if err != nil {
			_ = onExit()
		}
	}()

	if vp, err = impls.LoadYamlConfig(config, "Configuration"); err != nil {
		return err
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

	count := vp.GetInt("bp.count")
	if count <= 0 {
		return fmt.Errorf("invalid bp.count")
	}
	interval := vp.GetDuration("bp.interval")
	if interval <= 0 {
		return fmt.Errorf("invalid bp.interval")
	}
	if _Handler, err = NewHandler(count, interval); err != nil {
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
