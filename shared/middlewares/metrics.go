package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"tech-ip-sem2/shared/metrics"
)

func Metrics(next http.Handler, routeNameFunc func(r *http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := routeNameFunc(r)

		metrics.HTTPInFlight.Inc()
		defer metrics.HTTPInFlight.Dec()

		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		metrics.HTTPRequestsTotal.WithLabelValues(
			r.Method,
			route,
			strconv.Itoa(rw.statusCode),
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			r.Method,
			route,
		).Observe(duration)
	})
}
