package webserver

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// Cors respond an OPTIONS request (main auth)
func (s *WebServer) Cors(ctx *fasthttp.RequestCtx) {
	//return 200 with headers
	key, err := s.Auth(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.Response.Header.Set("X-RCP-Error", err.Error())
		return
	}

	for _, host := range key.AllowedHosts {
		if host == string(ctx.Request.Header.Peek("Origin")) {
			ctx.SetStatusCode(http.StatusNoContent)
			ctx.Response.Header.Set("Access-Control-Allow-Origin", host)
			return
		}
	}

	ctx.SetStatusCode(http.StatusForbidden)
	ctx.Response.Header.Set("X-RCP-Error", "host not allowed")
}
