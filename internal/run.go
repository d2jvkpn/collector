package internal

import (
	// "fmt"
	"context"
	"errors"
	"time"
)

func Run(addr string) (shutdown func() error, err error) {
	var shutdownHttp func() error

	if shutdownHttp, err = ServeHttp(addr); err != nil {
		return nil, err
	}

	if err = _KafkaHandler.Consume(); err != nil {
		return nil, errors.Join(err, shutdownHandler(), shutdownHttp())
	}

	shutdown = func() error {
		return errors.Join(shutdownHandler(), shutdownHttp())
	}

	return shutdown, nil
}

func shutdownHandler() (err error) {
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
