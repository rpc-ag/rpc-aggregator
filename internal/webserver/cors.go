package webserver

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// Cors respond an OPTIONS request (main auth)
func (s *WebServer) Cors(ctx *fasthttp.RequestCtx) {
	//return 200 with headers
	k := ctx.UserValue("api_key").(string)
	if k == "" {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.Response.Header.Set("X-RCP-Error", "no api key in request")
		return
	}

	key, found := s.auth.Auth(k)
	if !found {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.Response.Header.Set("X-RCP-Error", "api key not found")
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
