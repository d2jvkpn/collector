package kafka

import (
	"fmt"
	"log"
	"testing"

	"github.com/Shopify/sarama"
)

// go test -run TestProducer -- -addrs=localhost:29091
func TestProducer(t *testing.T) {
	var (
		err      error
		producer sarama.AsyncProducer
	)

	if producer, err = sarama.NewAsyncProducer(_TestAddrs, _TestConfig); err != nil {
		t.Fatal(err)
	}

	/*
		#### https://silverback-messaging.net/concepts/broker/kafka/kafka-partitioning.html?tabs=destination-partition-fluent%2Cenricher-fluent%2Cconcurrency-fluent%2Cassignment-fluent
		Kafka can guarantee ordering only inside the same partition and it is therefore important to
		be able to route correlated messages into the same partition. To do so you need to specify a
		key for each message and Kafka will put all messages with the same key in the same partition.
	*/
	for i := _TestIndex; i < _TestIndex+_TestNum; i++ {
		msg := fmt.Sprintf("hello message: %d", i)
		log.Println("--> send msg:", msg)

		pmsg1 := sarama.ProducerMessage{
			Topic: _TestTopic,
			Key:   sarama.StringEncoder("key0001"),
			Value: sarama.ByteEncoder([]byte(msg)),
		}

		producer.Input() <- &pmsg1

		// pmsg2 := pmsg1
		// pmsg2.Key = sarama.StringEncoder("key0002")
		// producer.Input() <- &pmsg2
	}

	if err = producer.Close(); err != nil {
		t.Fatal(err)
	}
}
