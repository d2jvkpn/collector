package internal

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

var (
	_TestTopic  string
	_TestKey    string
	_TestAddrs  []string
	_TestConfig *sarama.Config
)

func TestTimeFormat(t *testing.T) {
	now := time.Now()
	at := now.UTC()
	fmt.Println(">>>", at.Month())
	s := fmt.Sprintf("%dS%d", at.Year(), at.Month()%3)
	fmt.Println("   ", s)
}

func TestKafka(t *testing.T) {
	var (
		err      error
		producer sarama.AsyncProducer
	)

	_TestTopic = "collector"
	_TestKey = "key0001"
	_TestAddrs = []string{"localhost:29091"}
	_TestConfig = sarama.NewConfig()

	if _TestConfig.Version, err = sarama.ParseKafkaVersion("3.4.0"); err != nil {
		t.Fatal(err)
	}

	if producer, err = sarama.NewAsyncProducer(_TestAddrs, _TestConfig); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		data := NewData("test01", "biz0001").
			WithEventId(fmt.Sprintf("evnet%40d", i+1)).
			WithSvcV("0.1.0").
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
