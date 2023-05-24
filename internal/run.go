package internal

import (
	// "fmt"
	"context"
	"errors"
	"time"

	"github.com/d2jvkpn/collector/pkg/wrap"
)

func Run(addr string) (shutdown func() error, err error) {
	var shutdownProm func() error

	defer func() {
		if err != nil {
			_ = onExit()
		}
	}()

	if shutdownProm, err = wrap.PromHttp(addr); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = shutdownProm()
		}
	}()

	if err = _KafkaHandler.Consume(); err != nil {
		return nil, err
	}

	shutdown = func() error {
		return errors.Join(onExit(), shutdownProm())
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

	if _Logger != nil {
		e3 = _Logger.Down()
	}

	return errors.Join(e1, e2, e3)
}
