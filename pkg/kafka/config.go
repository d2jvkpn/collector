package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

type Config struct {
	Addrs   []string `mapstructure:"addrs"`
	Version string   `mapstructure:"version"` // 3.4.0
	Topic   string   `mapstructure:"topic"`

	// consumer
	GroupId string `mapstructure:"group_id"` // default
	// producer
	Key string `mapstructure:"key"`
}

func HandlerFromConfig(ctx context.Context, vp *viper.Viper, field string) (
	handler *Handler, err error) {
	var (
		config Config
		cfg    *sarama.Config
		group  sarama.ConsumerGroup
	)

	if err = vp.UnmarshalKey(field, &config); err != nil {
		return nil, err
	}

	if len(config.Addrs) == 0 || config.Version == "" {
		return nil, fmt.Errorf("invlaid addrs or version")
	}

	if config.Topic == "" {
		return nil, fmt.Errorf("invlaid topic")
	}

	cfg = sarama.NewConfig()
	if cfg.Version, err = sarama.ParseKafkaVersion(config.Version); err != nil {
		return nil, err
	}

	if group, err = sarama.NewConsumerGroup(config.Addrs, config.GroupId, cfg); err != nil {
		return nil, err
	}

	handler = NewHandler(ctx, group, []string{config.Topic})

	// alter handler.Logger later
	return handler, nil
}
