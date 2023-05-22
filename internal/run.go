package internal

import (
	// "fmt"
	"context"
	"errors"
	"time"
)

func Run() (err error) {
	if err = _KafkaHandler.Consume(); err != nil {
		return err
	}

	return nil
}

func Shutdown() (err error) {
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
