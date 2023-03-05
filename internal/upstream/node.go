package upstream

import (
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rpc-ag/rpc-aggregator/internal/config"
	"github.com/tufanbarisyildirim/balancer"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

var _ balancer.Node = (*Node)(nil)

// Node main node struct
type Node struct {
	Name         string `json:"name"`
	Chain        string `json:"chain"`
	Provider     string `json:"provider"`
	Endpoint     string `json:"endpoint"`
	Protocol     string `json:"protocol"`
	isHealthy    bool
	totalRequest uint64
	RateLimiter  *rate.Limiter
}

// ServeHTTP server http (actual proxy) through this node
// upstreamTook is only for the upstream request, not total (body reading & writing back to the client etc)
// it can be misleading when you got an error, so, count success only for accurate latency
func (n *Node) ServeHTTP(ctx *fasthttp.RequestCtx) (upstreamTook time.Duration, err error) {
	r := fasthttp.AcquireRequest()
	ctx.Request.CopyTo(r)
	r.SetRequestURI(n.Endpoint)
	r.SetTimeout(time.Second * 2)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release
	defer fasthttp.ReleaseRequest(r)
	start := time.Now()
	err = fasthttp.Do(r, resp)
	upstreamTook = time.Since(start)
	if err != nil {
		return upstreamTook, err
	}

	ctx.Response.Header.SetStatusCode(resp.StatusCode())
	resp.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.Add(string(key), string(value))
	})

	err = resp.BodyWriteTo(ctx.Response.BodyWriter())
	if err != nil {
		return upstreamTook, err
	}

	return upstreamTook, nil
}

// NewNode create new node
func NewNode(node *config.Node) (*Node, error) {
	n := &Node{
		Name:         node.Name,
		Chain:        node.Chain,
		Provider:     node.Provider,
		Endpoint:     node.Endpoint,
		Protocol:     node.Protocol,
		isHealthy:    true,
		totalRequest: 0,
		RateLimiter:  rate.NewLimiter(rate.Every(node.RateLimit.Per), node.RateLimit.Rate), //  ratelimit.New(node.RateLimit.Rate, ratelimit.Per(node.RateLimit.Per)),
	}

	return n, nil
}

// IsHealthy check if node can accept request
func (n *Node) IsHealthy() bool {
	return n.isHealthy
}

// TotalRequest total request done to this node so far
func (n *Node) TotalRequest() uint64 {
	return atomic.LoadUint64(&n.totalRequest)
}

// AverageResponseTime average response time of this node
func (n *Node) AverageResponseTime() time.Duration {
	return time.Millisecond * 200
}

// Load get load on that server
func (n *Node) Load() int64 {
	return 0
}

// NodeID a unique id for that particular node
func (n *Node) NodeID() string {
	return n.Name
}

// ProviderName node provider name
func (n *Node) ProviderName() string {
	return n.Provider
}

// SetHealthy set up/down the node.
func (n *Node) SetHealthy(healthy bool) {
	n.isHealthy = healthy
}

// HealthCheck a dummy healthcheck (it just recovers for now)
func (n *Node) HealthCheck() {
	//todo: do health check here, set SetHealthy(true) if pass
	n.SetHealthy(true)
}

// ToPromLabels return prometheus labels for the node
func (n *Node) ToPromLabels() prometheus.Labels {
	return prometheus.Labels{
		"chain":    n.Chain,
		"provider": n.Provider,
		"node_id":  n.NodeID(),
	}
}
