package wrap

import (
	// "fmt"
	"net"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type PromFasthttpOpt func(*fasthttp.Server)

func PromFasthttp(addr string, opts ...PromFasthttpOpt) (shutdown func() error, err error) {
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

	for i := range opts {
		opts[i](server)
	}

	shutdown = func() error {
		return server.Shutdown()
	}

	go func() {
		_ = server.Serve(listener)
	}()

	return shutdown, nil
}
