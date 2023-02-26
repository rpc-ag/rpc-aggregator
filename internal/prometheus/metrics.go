package prometheus

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	NodeRequests prometheus.Histogram
}

func NewMetrics() *Metrics {
	return &Metrics{
		NodeRequests: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "provider",
				Name:      "node_histogram",
				Help:      "Node latency",
			}),
	}
}
