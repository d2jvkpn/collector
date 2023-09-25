package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	// "github.com/d2jvkpn/collector/models"
	"github.com/d2jvkpn/collector/proto"

	"github.com/IBM/sarama"
)

func TestKafka(t *testing.T) {
	var (
		err      error
		producer sarama.AsyncProducer
	)

	if producer, err = sarama.NewAsyncProducer(_TestAddrs, _TestKafkaConfig); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		data := proto.NewRecordData("test01", "biz0001").
			WithEventId(fmt.Sprintf("event%04d", i+1)).
			WithSvcV("0.1.0").
			WithBizV("0.2.0").
			WithData(map[string]string{"hello": "world"})

		msg, _ := json.Marshal(data)

		pmsg1 := sarama.ProducerMessage{
			Topic: _TestTopic,
			Key:   sarama.StringEncoder(_TestKey),
			Value: sarama.ByteEncoder(msg),
		}

		producer.Input() <- &pmsg1
	}

	if err = producer.Close(); err != nil {
		t.Fatal(err)
	}
}
