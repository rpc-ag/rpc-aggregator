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
	auth     *config.Auth
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

func New(config *config.Config, auth *config.Auth, logger *zap.Logger) (*WebServer, error) {
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
		auth:     auth,
		server:   server,
		router:   r,
		upstream: &proxy.Upstream{Balancer: b},
	}
	ws.router.NotFound = ws.NotFound
	ws.router.OPTIONS("/{api_key}", ws.Cors)
	ws.router.ANY("/{api_key}", ws.Serve)

	return ws, nil
}

func applyFastHTTPConfig(server *fasthttp.Server, config *config.Config) {
	server.ReadTimeout = config.Webserver.ReadTimeout
	server.NoDefaultServerHeader = true
	if server.ReadTimeout == 0 {
		//TODO: move this to config
		server.ReadTimeout = 5 * time.Second
	}
}

// NotFound url handler, can't we handle all here?
func (s *WebServer) NotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

// Run starts server
func (s *WebServer) Run() error {
	return s.server.ListenAndServe(s.config.Webserver.Addr)
}

// Close sends `stop` signal to fasthttp server
func (s *WebServer) Close() error {
	return s.server.Shutdown()
}

func (s *WebServer) StartHealthChecker() {
	for {
		<-time.After(time.Second * 10) //todo: move this to config
		for _, n := range s.upstream.Balancer.UpstreamPool {
			if !n.IsHealthy() { //do check only if it is not healthy
				n.(*proxy.Node).HealthCheck()
				s.logger.Info("node is back", zap.String("node-id", n.NodeID()))
			}
		}
	}
}
