package metrics

type Metrics interface {
	IncrementRequestCount(method, path, status string)
	ObserveRequestDuration(method, path string, duration float64)
	IncrementActiveConnections(method, path string)
	DecrementActiveConnections(method, path string)
	ObserveRequestSize(method, path string, size float64)
	ObserveResponseSize(method, path string, size float64)
}
