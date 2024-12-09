package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTPRequestTotal Общее количество запросов
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HttpResponseDuration Время обработки запросов
	HTTPResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Histogram of response durations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func Register() {
	prometheus.MustRegister(HTTPRequestTotal)
	prometheus.MustRegister(HTTPResponseDuration)
}
