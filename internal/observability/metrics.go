package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DurationMetrics = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "function_duration_metric",
			Help:    "Function Execution Time",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"function", "status"},
	)
)

func init() {
	prometheus.MustRegister(DurationMetrics)
}

func GetDurationMetrics() *prometheus.HistogramVec {
	return DurationMetrics
}
