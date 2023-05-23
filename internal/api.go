package internal

import (
	// "fmt"
	// "log"
	"net"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func ServeHttp(addr string) (shutdown func() error, err error) {
	var (
		listener net.Listener
		server   *fasthttp.Server
	)

	if listener, err = net.Listen("tcp", addr); err != nil {
		return nil, err
	}

	server = new(fasthttp.Server)
	prom := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	server.Handler = fasthttp.CompressHandler(prom)

	shutdown = func() error {
		return server.Shutdown()
	}

	go func() {
		_ = server.Serve(listener)
	}()

	return shutdown, nil
}
