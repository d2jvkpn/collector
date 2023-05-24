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
		config *Config
		scfg   *sarama.Config
	)

	if config, scfg, err = NewConfigFromViper(vp, field); err != nil {
		return nil, err
	}

	if config.Key == "" {
		return nil, fmt.Errorf("invlaid topic or key")
	}

	producer = &KafkaProducer{config: *config}

	producer.producer, err = sarama.NewAsyncProducer(producer.config.Addrs, scfg)
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
