package proxy

import "github.com/tufanbarisyildirim/balancer"

// Upstream main upstream contains multiple upstreams in a balancer
type Upstream struct {
	Balancer *balancer.Balancer
}

//type Proxy func(ctx context.Context) error
//
//// proxy try next endpoint or all until it timed out
//func (b *Balancer) proxy(ctx context.Context, proxy Proxy) (dto, error) {
//
//	e := make(chan error)
//
//	go func() {
//		dtoResp, err := proxy(ctx)
//		if err != nil {
//			e <- err
//			return
//		}
//		c <- dtoResp
//	}()
//
//	select {
//	case <-ctx.Done():
//		return nil, ctx.Err()
//	case d := <-c:
//		return d, nil
//	case err := <-e:
//		return nil, err
//	}
//}
