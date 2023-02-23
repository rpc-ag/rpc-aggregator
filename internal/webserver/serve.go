package webserver

import (
	"github.com/rpc-ag/rpc-proxy/pkg/proxy"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Serve process the actual request
func (s *WebServer) Serve(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	next := s.upstream.Balancer.Next("a")
	if next == nil {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		s.logger.Error("no health node")
		return
	}
	node, ok := next.(*proxy.Node)

	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	err := node.ServeHTTP(ctx)
	if err != nil {
		node.SetHealthy(false)

		if err == fasthttp.ErrTimeout {
			s.Serve(ctx) // just try again with a new node
			return
		}

		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		s.logger.Error("failed serving", zap.Error(err))
		node.SetHealthy(false)
		return
	}
}
