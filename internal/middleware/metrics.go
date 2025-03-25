package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/leo-andrei/api-gateway/internal/logging"
	"github.com/leo-andrei/api-gateway/internal/metrics"
	"github.com/leo-andrei/api-gateway/pkg/responsewriter"
)

// MetricsMiddleware creates middleware for tracking metrics and logging
func MetricsMiddleware(next http.Handler, metrics metrics.Metrics, logger logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Track active connections
		path := r.URL.Path
		method := r.Method
		metrics.IncrementActiveConnections(method, path)
		defer metrics.DecrementActiveConnections(method, path)

		// Request size
		requestSize := 0
		if r.ContentLength > 0 {
			requestSize = int(r.ContentLength)
		}
		metrics.ObserveRequestSize(method, path, float64(requestSize))

		// Create a response writer wrapper to capture the status code and size
		rw := responsewriter.NewResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Request duration
		duration := time.Since(start)
		metrics.ObserveRequestDuration(method, path, duration.Seconds())

		// Request count
		status := fmt.Sprintf("%d", rw.StatusCode())
		metrics.IncrementRequestCount(method, path, status)

		// Response size
		metrics.ObserveResponseSize(method, path, float64(rw.Size()))

		// Log the request
		logger.LogRequest(r, duration, rw.StatusCode(), rw.Size())
	})
}
