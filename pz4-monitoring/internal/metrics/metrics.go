package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	HttpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_errors_total",
			Help: "Total number of HTTP error responses",
		},
		[]string{"method", "path", "status_code"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	StudentRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_student_requests_total",
			Help: "Total requests to student endpoint by student id",
		},
		[]string{"student_id"},
	)

	StudentHandlerDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_student_handler_duration_seconds",
			Help:    "Duration of student handler business logic",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"student_id"},
	)

	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_active_requests",
			Help: "Number of HTTP requests currently being processed",
		},
	)
)
