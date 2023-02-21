package proxy

import (
	"sync/atomic"
	"time"

	"github.com/tufanbarisyildirim/balancer"
	"github.com/valyala/fasthttp"
)

var _ balancer.Node = (*Node)(nil)

type Node struct {
	Name         string `json:"name"`
	Chain        string `json:"chain"`
	Provider     string `json:"provider"`
	Endpoint     string `json:"endpoint"`
	Protocol     string `json:"protocol"`
	isHealthy    bool
	totalRequest uint64
}

func (n *Node) ServeHTTP(ctx *fasthttp.RequestCtx) error {
	r := fasthttp.AcquireRequest()
	ctx.Request.CopyTo(r)
	r.SetRequestURI(n.Endpoint)
	r.SetTimeout(time.Second * 2)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release
	defer fasthttp.ReleaseRequest(r)
	err := fasthttp.Do(r, resp)

	if err != nil {
		return err
	}

	ctx.Response.Header.SetStatusCode(resp.StatusCode())
	resp.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.Add(string(key), string(value))
	})

	err = resp.BodyWriteTo(ctx.Response.BodyWriter())
	if err != nil {
		return err
	}

	return nil
}

func NewNode(name, chain, provider, endpoint, protocol string) (*Node, error) {
	return &Node{
		Name:         name,
		Chain:        chain,
		Provider:     provider,
		Endpoint:     endpoint,
		Protocol:     protocol,
		isHealthy:    true,
		totalRequest: 0,
	}, nil
}

func (n *Node) IsHealthy() bool {
	return n.isHealthy
}

func (n *Node) TotalRequest() uint64 {
	return atomic.LoadUint64(&n.totalRequest)
}

func (n *Node) AverageResponseTime() time.Duration {
	return time.Millisecond * 200
}

func (n *Node) Load() int64 {
	return 0
}

func (n *Node) NodeID() string {
	return n.Name
}

func (n *Node) ProviderName() string {
	return n.Provider
}

func (n *Node) SetHealthy(healthy bool) {
	n.isHealthy = healthy
}

func (n *Node) HealthCheck() {
	//todo: do health check here, set SetHealthy(true) if pass
	n.SetHealthy(true)
}
