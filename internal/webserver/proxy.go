package webserver

import (
	"errors"
	"net/http"

	"github.com/rpc-ag/rpc-proxy/internal/config"
	"github.com/rpc-ag/rpc-proxy/pkg/upstream"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Proxy process the actual request
func (s *WebServer) Proxy(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	next := s.upstream.Balancer.Next("a")
	if next == nil {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		s.logger.Error("no health node")
		return
	}
	node, ok := next.(*upstream.Node)

	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	key, err := s.Auth(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.Response.Header.Set("X-RCP-Error", err.Error())
		return
	}

	// TODO: if a node is not accepting a new request due to rate limiting
	// decide to bypass or just remove from the list?
	// can healtchchecker check limiter and add back if allowed?
	// for now, just remove the node from balancer list and try the next one
	// OH! we do auth again, separate request/auth and do auth only once
	if !node.RateLimiter.Allow() {
		node.SetHealthy(false)
		s.Proxy(ctx)
		return
	}

	if !key.RateLimiter.Allow() {
		ctx.SetStatusCode(http.StatusTooManyRequests)
		ctx.Response.Header.Set("X-RCP-Error", "rate limit exceeded")
		return
	}

	err = node.ServeHTTP(ctx)
	if err != nil {
		node.SetHealthy(false)

		if err == fasthttp.ErrTimeout {
			s.Proxy(ctx) // just try again with a new node
			return
		}

		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		s.logger.Error("failed serving", zap.Error(err))
		node.SetHealthy(false)
		return
	}
}

// Auth check authentication first
// TODO: can that be a lightweight middleware please?
func (s *WebServer) Auth(ctx *fasthttp.RequestCtx) (*config.APIKey, error) {
	k := ctx.UserValue("api_key").(string)
	if k == "" {
		return nil, errors.New("no ap key in request")
	}

	key, found := s.auth.Auth(k)
	if !found {
		return nil, errors.New("api key not found")
	}

	return key, nil
}
