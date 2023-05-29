package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/d2jvkpn/collector/internal/settings"
	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/spf13/viper"
)

func Run(addr string) (shutdown func() error, err error) {
	var (
		vp              *viper.Viper
		shutdownMetrics func() error
	)

	defer func() {
		if err != nil {
			_ = onExit()
		}
	}()

	// shutdownMetrics, err = wrap.PromFasthttp(addr)
	if vp = settings.ConfigSub("metrics"); vp == nil {
		return nil, fmt.Errorf("config.metrics is unset")
	}

	vp.Set("addr", addr)
	if shutdownMetrics, err = wrap.HttpMetrics(vp, settings.Meta); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = shutdownMetrics()
		}
	}()

	_KafkaHandler.Consume()

	shutdown = func() error {
		return errors.Join(onExit(), shutdownMetrics())
	}

	return shutdown, nil
}

func onExit() (err error) {
	var e1, e2, e3 error

	if _KafkaHandler != nil {
		e1 = _KafkaHandler.Close()
	}

	if _Handler != nil {
		_Handler.Down()
	}

	if _MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		e2 = _MongoClient.Disconnect(ctx)
	}

	if settings.Logger != nil {
		e3 = settings.Logger.Down()
	}

	return errors.Join(e1, e2, e3)
}
