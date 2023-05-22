package wrap

import (
	"context"
	// "fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConf struct {
	Uri             string `mapstructure:"uri"`
	MinPoolSize     uint64 `mapstructure:"min_pool_size"`
	MaxPoolSize     uint64 `mapstructure:"max_pool_size"`
	TimeoutSecs     int    `mapstructure:"timeout_secs"`
	MaxConnIdleSecs int    `mapstructure:"max_conn_idle_secs"`
}

func MongoClient(vp *viper.Viper, field string) (client *mongo.Client, err error) {
	var (
		ctx    context.Context
		config MongoConf
		cancel func()
	)

	if err = vp.UnmarshalKey(field, &config); err != nil {
		return nil, err
	}

	timeout := time.Duration(config.TimeoutSecs) * time.Second
	maxConnIdle := time.Duration(config.MaxConnIdleSecs) * time.Second

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err = mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(config.Uri).
			SetServerSelectionTimeout(timeout).
			SetMinPoolSize(config.MinPoolSize).
			SetMaxPoolSize(config.MaxPoolSize).
			SetMaxConnIdleTime(maxConnIdle),
	)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
