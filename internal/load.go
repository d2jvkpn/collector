package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/d2jvkpn/collector/internal/biz"
	"github.com/d2jvkpn/collector/internal/settings"
	"github.com/d2jvkpn/collector/pkg/kafka"
	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/cloud-logging"
	"github.com/d2jvkpn/gotk/cloud-tracing"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadLocal(config string) (err error) {
	var vp *viper.Viper

	if vp, err = gotk.LoadYamlConfig(config, "Configuration"); err != nil {
		return err
	}

	return load(vp)
}

func LoadConsul(config string) (err error) {
	return fmt.Errorf("unimplemented")
}

func load(vp *viper.Viper) (err error) {
	defer func() {
		if err == nil {
			return
		}

		_ = onExit()
	}()

	vp.SetDefault("http.cors", "*")
	vp.SetDefault("log.size_mb", 256)

	if err = settings.SetConfig(vp); err != nil {
		return err
	}
	settings.Meta["service_name"] = settings.ServiceName()
	settings.Meta["startup"] = time.Now().Format(time.RFC3339)

	// if _ServiceName = _Config.GetString("service_name"); _ServiceName == "" {
	// 	return fmt.Errorf("service_name is empty in config")
	// }

	//
	settings.Logger, err = logging.NewLogger(
		vp.GetString("log.path"),
		zap.InfoLevel,
		vp.GetInt("log.size_mb"),
	)
	if err != nil {
		return fmt.Errorf("NewLogger: %w", err)
	}
	_Logger = settings.Logger.Named("internal")

	//
	if _MongoClient, err = wrap.MongoClient(vp, "mongodb"); err != nil {
		return fmt.Errorf("MongoClient: %w", err)
	}

	//
	if _RecordHandler, err = biz.NewRecordHandler(vp.Sub("bp")); err != nil {
		return fmt.Errorf("NewHandler: %w", err)
	}

	db := vp.GetString("mongodb.db")
	if db == "" {
		return fmt.Errorf("mongodb.db is empty")
	}

	database := _MongoClient.Database(db)
	err = _RecordHandler.
		WithLogger(settings.Logger.Named("record_handler")).
		WithDatabase(database).
		Ok()
	if err != nil {
		return err
	}

	//
	if _KafkaHandler, err = kafka.HandlerFromConfig(context.TODO(), vp, "kafka"); err != nil {
		return fmt.Errorf("HandlerFromConfig: %w", err)
	}

	err = _KafkaHandler.
		WithLogger(settings.Logger.Named("kafka_handler")).
		WithHandle(_RecordHandler.Handle).Ok()
	if err != nil {
		return fmt.Errorf("Handler: %w", err)
	}

	if _Metrics = settings.ConfigSub("metrics"); _Metrics == nil {
		return fmt.Errorf("config.metrics is unset")
	}
	settings.Meta["metrics_addr"] = _Metrics.GetString("addr")
	settings.Meta["metrics_prometheus"] = _Metrics.GetBool("prometheus")
	settings.Meta["metrics_debug"] = _Metrics.GetBool("debug")

	//
	otelConfig := settings.ConfigSub("otel")
	if otelConfig == nil {
		return fmt.Errorf("config.otel is unset")
	}
	if otelConfig.GetBool("enable") {
		settings.Meta["otel_enable"] = true
		otelAddr := otelConfig.GetString("addr")
		settings.Meta["otel_addr"] = otelAddr
		_CloseOtel, err = tracing.LoadOtelGrpc(otelAddr, settings.ServiceName(), false)
		if err != nil {
			return fmt.Errorf("LoadOtel: %w", err)
		}
	} else {
		settings.Meta["otel_enable"] = false
	}

	//
	grpcConfig := settings.ConfigSub("grpc")
	if grpcConfig == nil {
		return fmt.Errorf("config.grpc is unset")
	}
	_GrpcSS, err = biz.NewGSS(
		settings.Logger.Named("grpc_interceptor"),
		database,
		grpcConfig,
		otelConfig.GetBool("enable"),
	)
	if err != nil {
		return fmt.Errorf("NewGSS: %w", err)
	}

	return nil
}
