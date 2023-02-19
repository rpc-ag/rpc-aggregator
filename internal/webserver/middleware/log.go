package middleware

import (
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// NewLogger logs after a request
func NewLogger(logger *zap.Logger) Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			var startTime = time.Now()
			next(ctx)
			logger.Debug("access",
				zap.Int("code", ctx.Response.StatusCode()),
				zap.Duration("time", time.Since(startTime)),
				zap.ByteString("method", ctx.Method()),
				zap.ByteString("path", ctx.Path()))
		}
	}
}
