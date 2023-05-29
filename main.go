package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/d2jvkpn/collector/internal"
	"github.com/d2jvkpn/collector/internal/settings"

	"github.com/d2jvkpn/gotk"
)

var (
	//go:embed project.yaml
	_Project []byte
)

func init() {
	gotk.RegisterLogPrinter()
}

func main() {
	var (
		consul   bool
		config   string
		addr     string
		err      error
		quit     chan os.Signal
		shutdown func() error
	)

	if err = settings.SetProject(_Project); err != nil {
		log.Fatalln(err)
	}

	flag.StringVar(&config, "config", "configs/local.yaml", "configuration file path")
	flag.StringVar(&addr, "addr", "0.0.0.0:5011", "prometheus metrics http server")
	flag.BoolVar(&consul, "consul", false, "using consul")

	flag.Usage = func() {
		output := flag.CommandLine.Output()

		fmt.Fprintf(output, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(output, "\nConfiguration:\n```yaml\n%s```\n", settings.DemoConfig())
		fmt.Fprintf(output, "\nBuild:\n```text\n%s\n```\n", gotk.BuildInfoText(settings.Meta))
	}

	flag.Parse()

	settings.Meta["config"] = config
	settings.Meta["consul"] = consul

	if consul {
		err = internal.LoadConsul(config)
	} else {
		err = internal.LoadLocal(config)
	}
	if err != nil {
		log.Fatalln(err)
	}

	if shutdown, err = internal.Run(addr); err != nil {
		log.Fatalln(err)
	}
	log.Printf(">>> The server is starting, http listening on %s...\n", addr)

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)

	select {
	case sig := <-quit:
		// if sig == syscall.SIGUSR2 {...}
		fmt.Println("... received:", sig)
		err = shutdown()
	}

	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("<<< Exit")
	}
}
