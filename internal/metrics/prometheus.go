package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsService handles metrics collection and reporting
type MetricsService struct {
	RequestCount      *prometheus.CounterVec
	RequestDuration   *prometheus.HistogramVec
	RequestSize       *prometheus.SummaryVec
	ResponseSize      *prometheus.SummaryVec
	ActiveConnections *prometheus.GaugeVec
}

var _ Metrics = (*MetricsService)(nil)

// NewMetricsService initializes a new metrics service
func NewMetricsService() *MetricsService {
	m := &MetricsService{
		RequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_gateway_requests_total",
				Help: "Total number of requests processed by the API Gateway",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_gateway_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		RequestSize: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "api_gateway_request_size_bytes",
				Help:       "Request size in bytes",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"method", "path"},
		),
		ResponseSize: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "api_gateway_response_size_bytes",
				Help:       "Response size in bytes",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"method", "path"},
		),
		ActiveConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "api_gateway_active_connections",
				Help: "Number of active connections",
			},
			[]string{"method", "path"},
		),
	}

	// Register metrics with Prometheus
	prometheus.MustRegister(m.RequestCount)
	prometheus.MustRegister(m.RequestDuration)
	prometheus.MustRegister(m.RequestSize)
	prometheus.MustRegister(m.ResponseSize)
	prometheus.MustRegister(m.ActiveConnections)

	return m
}

func (m *MetricsService) IncrementRequestCount(method, path, status string) {
	m.RequestCount.WithLabelValues(method, path, status).Inc()
}

func (m *MetricsService) ObserveRequestDuration(method, path string, duration float64) {
	m.RequestDuration.WithLabelValues(method, path).Observe(duration)
}

func (m *MetricsService) IncrementActiveConnections(method, path string) {
	m.ActiveConnections.WithLabelValues(method, path).Inc()
}

func (m *MetricsService) DecrementActiveConnections(method, path string) {
	m.ActiveConnections.WithLabelValues(method, path).Dec()
}

func (m *MetricsService) ObserveRequestSize(method, path string, size float64) {
	m.RequestSize.WithLabelValues(method, path).Observe(size)
}

func (m *MetricsService) ObserveResponseSize(method, path string, size float64) {
	m.ResponseSize.WithLabelValues(method, path).Observe(size)
}
