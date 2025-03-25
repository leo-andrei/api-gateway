package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockMetrics is a mock implementation of the Metrics interface
type MockMetrics struct {
	mock.Mock
}

func (m *MockMetrics) IncrementRequestCount(method, path, status string) {
	m.Called(method, path, status)
}

func (m *MockMetrics) ObserveRequestDuration(method, path string, duration float64) {
	m.Called(method, path, duration)
}

func (m *MockMetrics) IncrementActiveConnections(method, path string) {
	m.Called(method, path)
}

func (m *MockMetrics) DecrementActiveConnections(method, path string) {
	m.Called(method, path)
}

func (m *MockMetrics) ObserveRequestSize(method, path string, size float64) {
	m.Called(method, path, size)
}

func (m *MockMetrics) ObserveResponseSize(method, path string, size float64) {
	m.Called(method, path, size)
}

// MockLogger is a mock implementation of the Logger interface
type MockLogger struct {
	mock.Mock
}

func (l *MockLogger) LogRequest(r *http.Request, duration time.Duration, status int, responseSize int) {
	l.Called(r, duration, status, responseSize)
}

func (l *MockLogger) Fatal(v string) {
	l.Called(v)
}

func (l *MockLogger) Fatalf(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *MockLogger) Info(v string) {
	l.Called(v)
}

func (l *MockLogger) Infof(format string, args ...interface{}) {
	l.Called(append([]interface{}{format}, args...)...)
}

func (l *MockLogger) Shutdown() {
	l.Called()
}

func TestMetricsMiddleware(t *testing.T) {
	mockMetrics := new(MockMetrics)
	mockLogger := new(MockLogger)

	// Set up expectations
	mockMetrics.On("IncrementActiveConnections", "GET", "/test").Once()
	mockMetrics.On("DecrementActiveConnections", "GET", "/test").Once()
	mockMetrics.On("IncrementRequestCount", "GET", "/test", "200").Once()
	mockMetrics.On("ObserveRequestDuration", "GET", "/test", mock.Anything).Once()
	mockMetrics.On("ObserveRequestSize", "GET", "/test", mock.Anything).Once()
	mockMetrics.On("ObserveResponseSize", "GET", "/test", mock.Anything).Once()
	mockLogger.On("LogRequest", mock.Anything, mock.Anything, 200, 0).Once()

	// Create a test handler
	handler := MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), mockMetrics, mockLogger)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Assert expectations
	mockMetrics.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
