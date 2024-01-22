package middleware

import "github.com/valyala/fasthttp"

// Middleware chainable fasthttp middleware func
type Middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler

// Apply process multiple middleware
func Apply(h fasthttp.RequestHandler, m ...Middleware) fasthttp.RequestHandler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

// Onion Onion style middleware
func Onion(m ...Middleware) Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return Apply(next, m...)
	}
}
