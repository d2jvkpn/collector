package internal

import (
	"context"
	"fmt"
	// "time"

	"github.com/d2jvkpn/collector/internal/biz"
	"github.com/d2jvkpn/collector/internal/settings"
	"github.com/d2jvkpn/collector/pkg/kafka"
	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/cloud-logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadLocal(config, addr string) (err error) {
	var vp *viper.Viper

	if vp, err = gotk.LoadYamlConfig(config, "Configuration"); err != nil {
		return err
	}

	return load(vp, addr)
}

func LoadConsul(config, addr string) (err error) {
	return fmt.Errorf("unimplemented")
}

func load(vp *viper.Viper, addr string) (err error) {
	defer func() {
		if err != nil {
			_ = onExit()
		}
	}()

	vp.SetDefault("http.cors", "*")
	vp.SetDefault("log.size_mb", 256)

	if err = settings.SetConfig(vp); err != nil {
		return err
	}

	// if _ServiceName = _Config.GetString("service_name"); _ServiceName == "" {
	// 	return fmt.Errorf("service_name is empty in config")
	// }

	settings.Logger, err = logging.NewLogger(
		vp.GetString("log.path"),
		zap.InfoLevel,
		vp.GetInt("log.size_mb"),
	)
	if err != nil {
		return fmt.Errorf("NewLogger: %w", err)
	}
	_Logger = settings.Logger.Named("internal")

	if _MongoClient, err = wrap.MongoClient(vp, "mongodb"); err != nil {
		return fmt.Errorf("MongoClient: %w", err)
	}

	//
	if _Handler, err = biz.NewHandler(vp.Sub("bp")); err != nil {
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

	err = _KafkaHandler.WithLogger(_Logger.Named("kafka")).WithHandle(_Handler.Handle).Ok()
	if err != nil {
		return fmt.Errorf("Handler: %w", err)
	}

	if _Metrics = settings.ConfigSub("metrics"); vp == nil {
		return fmt.Errorf("config.metrics is unset")
	}
	_Metrics.Set("addr", addr)
	settings.Meta["metrics_addr"] = _Metrics.GetString("addr")
	settings.Meta["metrics_prometheus"] = _Metrics.GetBool("prometheus")
	settings.Meta["metrics_debug"] = _Metrics.GetBool("debug")

	settings.Meta["service_name"] = settings.ServiceName()

	return nil
}
