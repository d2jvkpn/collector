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
	"github.com/d2jvkpn/collector/pkg/wrap"

	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
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
		config   string
		addr     string
		err      error
		meta     map[string]any
		quit     chan os.Signal
		shutdown func() error
		project  *viper.Viper
	)

	if project, err = wrap.LoadYamlBytes(_Project); err != nil {
		log.Fatalln(err)
	}

	meta = gotk.BuildInfo()
	meta["project"] = project.GetString("project")
	meta["version"] = project.GetString("version")

	flag.StringVar(&config, "config", "configs/local.yaml", "configuration file path")
	flag.StringVar(&addr, "addr", "0.0.0.0:5011", "prometheus metrics http server")

	flag.Usage = func() {
		output := flag.CommandLine.Output()

		fmt.Fprintf(output, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(output, "\nConfiguration:\n```yaml\n%s```\n", project.GetString("config"))
		fmt.Fprintf(output, "\nBuild:\n```text\n%s\n```\n", gotk.BuildInfoText(meta))
	}

	flag.Parse()

	if err = internal.Load(config); err != nil {
		log.Fatalln(err)
	}

	if shutdown, err = internal.Run(addr); err != nil {
		log.Fatalln(err)
	}
	log.Printf(">>> The server is starting, http listening on %s...\n", addr)

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)

	select {
	case <-quit:
		fmt.Println("...")
		err = shutdown()
	}

	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("<<< Exit")
	}
}
