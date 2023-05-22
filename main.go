package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/d2jvkpn/collector/internal"
	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

var (
	//go:embed project.yaml
	_ProjectBts []byte
	_Project    *viper.Viper
)

func init() {
	gotk.RegisterLogPrinter()

	_Project = viper.New()
	_Project.SetConfigType("yaml")
}

func main() {
	var (
		config string
		err    error
		quit   chan os.Signal
	)

	if err = _Project.ReadConfig(bytes.NewReader(_ProjectBts)); err != nil {
		log.Fatalln(err)
	}

	flag.StringVar(&config, "config", "configs/local.yaml", "configuration file path")

	flag.Usage = func() {
		output := flag.CommandLine.Output()

		fmt.Fprintf(output, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(output, "\nConfiguration:\n```yaml\n%s```\n", _Project.GetString("config"))
	}

	flag.Parse()

	if err = internal.Load(config); err != nil {
		_ = internal.Shutdown()
		log.Fatalln(err)
	}

	if err = internal.Run(); err != nil {
		_ = internal.Shutdown()
		log.Fatalln(err)
	}
	log.Println(">>> The server is starting...")

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)

	select {
	case <-quit:
		fmt.Println("...")
		err = internal.Shutdown()
		log.Printf("<<< Exit\n")
	}

	if err != nil {
		log.Fatalln(err)
	}
}
