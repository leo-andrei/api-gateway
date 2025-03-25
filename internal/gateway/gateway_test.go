package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/leo-andrei/api-gateway/config"
	"github.com/leo-andrei/api-gateway/internal/logging"
	"github.com/leo-andrei/api-gateway/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestGatewayHealthEndpoint(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Server: struct {
			Port int `yaml:"port"`
		}{
			Port: 8080,
		},
		Logging: struct {
			Level  string `yaml:"level"`
			Format string `yaml:"format"`
		}{
			Level:  "info",
			Format: "json",
		},
		Routes: []config.Route{},
	}

	logService := logging.NewLogService(cfg.Logging)
	metricsService := metrics.NewMetricsService()
	gw := NewGateway(cfg, logService, metricsService)
	gw.SetupRoutes()

	// Create a request to the health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	gw.router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestNewGateway(t *testing.T) {
	cfg := &config.Config{
		Server: struct {
			Port int `yaml:"port"`
		}{
			Port: 8080,
		},
		Logging: struct {
			Level  string `yaml:"level"`
			Format string `yaml:"format"`
		}{
			Level:  "info",
			Format: "json",
		},
		Routes: []config.Route{},
	}
	logService := logging.NewLogService(cfg.Logging)
	metricsService := metrics.NewMetricsService()

	gw := NewGateway(cfg, logService, metricsService)

	assert.NotNil(t, gw)
	assert.Equal(t, cfg, gw.config)
	assert.NotNil(t, gw.router)
	assert.NotNil(t, gw.logService)
	assert.NotNil(t, gw.metricsService)
}

func TestGateway_Shutdown(t *testing.T) {
	cfg := &config.Config{
		Server: struct {
			Port int `yaml:"port"`
		}{
			Port: 8080,
		},
		Logging: struct {
			Level  string `yaml:"level"`
			Format string `yaml:"format"`
		}{
			Level:  "info",
			Format: "json",
		},
		Routes: []config.Route{},
	}
	logService := logging.NewLogService(cfg.Logging)
	metricsService := metrics.NewMetricsService()
	gw := NewGateway(cfg, logService, metricsService)

	gw.server = &http.Server{
		Addr:    ":8080",
		Handler: gw.router,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := gw.Shutdown(ctx)
	assert.NoError(t, err)
}
