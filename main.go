package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/d2jvkpn/collector/internal"
)

func main() {
	var (
		config string
		err    error
		quit   chan os.Signal
	)

	flag.StringVar(&config, "config", "configs/local.yaml", "configuration file path")
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
