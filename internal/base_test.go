package internal

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/Shopify/sarama"
	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	_TestTopic       string
	_TestKey         string
	_TestAddrs       []string
	_TestFlag        *flag.FlagSet
	_TestCtx         context.Context
	_TestConfig      *viper.Viper
	_TestKafkaConfig *sarama.Config
)

func TestFormat(t *testing.T) {
	now := time.Now()
	at := now.UTC()

	fmt.Println(">>>", at.Month())
	fmt.Println("   ", fmt.Sprintf("%dS%d", at.Year(), at.Month()%3))

	fmt.Printf("~~~ %04d\n", 3)
}

func TestMain(m *testing.M) {
	_TestTopic = "collector"
	_TestKey = "key0001"
	_TestAddrs = []string{"localhost:29091"}
	_TestCtx = context.Background()

	var (
		configFile string
		err        error
	)

	defer func() {
		if err != nil {
			fmt.Printf("!!! TestMain: %v\n", err)
			os.Exit(1)
		}
	}()

	_TestFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	_TestFlag.StringVar(&configFile, "config", "configs/local.yaml", "config filepath")
	if configFile, err = gotk.RootFile(configFile); err != nil {
		return
	}

	if _TestConfig, err = wrap.LoadYamlConfig(configFile, "TestConfig"); err != nil {
		return
	}

	_TestFlag.Parse(flag.Args())

	_TestKafkaConfig = sarama.NewConfig()
	_TestKafkaConfig.Version, err = sarama.ParseKafkaVersion(_TestConfig.GetString("kafka.version"))
	if err != nil {
		return
	}

	fmt.Printf("~~~ config %s\n", configFile)

	m.Run()
}
