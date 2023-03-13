package prometheus

import "github.com/prometheus/client_golang/prometheus"

// Metrics all prometheus metrics
type Metrics struct {
	NodeRequests *prometheus.HistogramVec
}

// NewMetrics Creates new metrics holder
func NewMetrics() *Metrics {
	return &Metrics{
		NodeRequests: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "provider",
				Name:      "request_duration_seconds",
				Help:      "Node request duration",
			}, []string{"chain", "provider", "node_id"}),
	}
}
