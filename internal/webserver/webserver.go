package webserver

import (
	"fmt"
	"time"

	"github.com/fasthttp/router"
	"github.com/rpc-ag/rpc-proxy/internal/config"
	"github.com/rpc-ag/rpc-proxy/internal/webserver/middleware"
	"github.com/rpc-ag/rpc-proxy/pkg/proxy"
	"github.com/tufanbarisyildirim/balancer"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type WebServer struct {
	logger   *zap.Logger
	config   *config.Config
	server   *fasthttp.Server
	router   *router.Router
	upstream *proxy.Upstream
}

type loggerAdapter struct {
	logger *zap.Logger
}

func (l loggerAdapter) Printf(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func New(config *config.Config, logger *zap.Logger) (*WebServer, error) {
	server := new(fasthttp.Server)
	applyFastHTTPConfig(server, config)
	server.Logger = loggerAdapter{logger: logger}

	r := router.New()

	var middlewares = middleware.Onion(
		middleware.NewLogger(logger),
		middleware.NewRecover(logger),
	)

	server.Handler = middlewares(r.Handler)

	b := balancer.NewBalancer()
	for _, n := range config.Nodes {
		node, err := proxy.NewNode(n.Name, n.Chain, n.Provider, n.Endpoint, n.Protocol)
		if err != nil {
			logger.Error("error creating node", zap.Any("node", node), zap.Error(err))
			continue
		}
		b.Add(node)
	}

	ws := &WebServer{
		logger:   logger,
		config:   config,
		server:   server,
		router:   r,
		upstream: &proxy.Upstream{Balancer: b},
	}
	ws.router.NotFound = ws.NotFound

	return ws, nil
}

func applyFastHTTPConfig(server *fasthttp.Server, config *config.Config) {
	server.ReadTimeout = config.Webserver.ReadTimeout
	server.NoDefaultServerHeader = true
	if server.ReadTimeout == 0 {
		server.ReadTimeout = 5 * time.Second
	}
}

// NotFound url handler, can't we handle all here?
func (s *WebServer) NotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	next := s.upstream.Balancer.Next("a")
	if next == nil {
		ctx.SetBody([]byte("no upstream found"))
		return
	}
	if node, ok := next.(*proxy.Node); ok {
		node.ServeHTTP(ctx)
		//ctx.SetBody([]byte(fmt.Sprintf("upstream %s will be called", node.Endpoint)))
	} else {
		ctx.SetBody([]byte("invalid upstream"))
	}

}

// Run starts server
func (s *WebServer) Run() error {
	return s.server.ListenAndServe(s.config.Webserver.Addr)
}

// Close sends `stop` signal to fasthttp server
func (s *WebServer) Close() error {
	return s.server.Shutdown()
}
