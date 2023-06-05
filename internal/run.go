package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/d2jvkpn/collector/internal/settings"

	"github.com/d2jvkpn/gotk/cloud-metrics"
)

func Run(addr string) (shutdown func() error, err error) {
	var (
		listener        net.Listener
		shutdownMetrics func() error
	)

	defer func() {
		if err == nil {
			return
		}

		if shutdownMetrics != nil {
			_ = shutdownMetrics()
		}
		_ = onExit()
	}()

	// shutdownMetrics, err = wrap.PromFasthttp(addr)
	if shutdownMetrics, err = metrics.HttpMetrics(_Metrics, settings.Meta); err != nil {
		return nil, err
	}

	_KafkaHandler.Consume()

	if listener, err = net.Listen("tcp", addr); err != nil {
		return nil, fmt.Errorf("net.Listen: %w", err)
	}

	go func() {
		_ = _GrpcSS.Serve(listener)
	}()

	_Logger.Info(fmt.Sprintf("GRPC server is listening on %s", addr))
	shutdown = func() error {
		_GrpcSS.Shutdown()
		e1 := shutdownMetrics()
		e2 := onExit()

		return errors.Join(e1, e2)
	}

	return shutdown, nil
}

func onExit() (err error) {
	_Logger.Info(fmt.Sprintf("GRPC server is shutdowning donw"))

	if _KafkaHandler != nil {
		err = errors.Join(err, _KafkaHandler.Close())
	}

	if _RecordHandler != nil {
		_RecordHandler.Down()
	}

	if _MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		err = errors.Join(err, _MongoClient.Disconnect(ctx))
	}

	if _CloseOtel != nil {
		err = errors.Join(err, _CloseOtel())
	}

	if settings.Logger != nil {
		err = errors.Join(err, settings.Logger.Down())
	}

	return err
}
