package middleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// NewRecover recover middleware, prevent crashing when a handler paniced
func NewRecover(logger *zap.Logger) Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			defer func() {
				err := recover()
				if err == nil {
					return
				}

				if err, converted := err.(error); converted {
					logger.Error("recovered from panic", zap.Error(err),
						zap.Stack("stack"))
					return
				}

				logger.Error("panic on non-error", zap.Any("err", err))
			}()

			next(ctx)
		}
	}
}
