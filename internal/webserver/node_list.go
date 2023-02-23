package webserver

import (
	"encoding/json"

	"github.com/rpc-ag/rpc-proxy/pkg/upstream"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type dtoNode struct {
	Name         string `json:"name"`
	Chain        string `json:"chain"`
	Provider     string `json:"provider"`
	Protocol     string `json:"protocol"`
	IsHealthy    bool   `json:"is_healthy"`
	TotalRequest uint64 `json:"total_request"`
}

// NodeList list the nodes and general availability
func (s *WebServer) NodeList(ctx *fasthttp.RequestCtx) {
	var nodes []dtoNode
	for _, n := range s.upstream.Balancer.UpstreamPool {
		node := n.(*upstream.Node)
		nodes = append(nodes, dtoNode{
			Name:         node.Name,
			Chain:        node.Chain,
			Provider:     node.Provider,
			Protocol:     node.Protocol,
			IsHealthy:    node.IsHealthy(),
			TotalRequest: node.TotalRequest(),
		})
	}
	nodesJSON, err := json.Marshal(nodes)
	if err != nil {
		s.logger.Error("failed serializing nodes", zap.Error(err))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = ctx.Write(nodesJSON)
	if err != nil {
		s.logger.Error("failed writing node list response", zap.Error(err))
		return
	}
}
