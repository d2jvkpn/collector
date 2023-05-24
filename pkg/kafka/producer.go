package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

type KafkaProducer struct {
	config   Config
	producer sarama.AsyncProducer
}

func NewKafkaProducer(vp *viper.Viper, field string) (producer *KafkaProducer, err error) {
	var (
		config Config
		cfg    *sarama.Config
	)

	if err = vp.UnmarshalKey(field, &config); err != nil {
		return nil, err
	}

	if len(config.Addrs) == 0 || config.Version == "" {
		return nil, fmt.Errorf("invlaid addrs or version")
	}

	cfg = sarama.NewConfig()
	if cfg.Version, err = sarama.ParseKafkaVersion(config.Version); err != nil {
		return nil, err
	}

	producer = &KafkaProducer{config: config}
	if producer.config.Topic == "" || producer.config.Key == "" {
		return nil, fmt.Errorf("invlaid topic or key")
	}

	producer.producer, err = sarama.NewAsyncProducer(producer.config.Addrs, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (producer *KafkaProducer) SendMsg(bts []byte) (msg *sarama.ProducerMessage, ok bool) {
	msg = &sarama.ProducerMessage{
		Topic: producer.config.Topic,
		Key:   sarama.StringEncoder(producer.config.Key),
		Value: sarama.ByteEncoder(bts),
	}

	producer.producer.Input() <- msg
	return msg, true
}

func (producer *KafkaProducer) Close() (err error) {
	return producer.producer.Close()
}
