package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricsService(t *testing.T) {
	metricsService := NewMetricsService()

	assert.NotNil(t, metricsService.RequestCount)
	assert.NotNil(t, metricsService.RequestDuration)
	assert.NotNil(t, metricsService.RequestSize)
	assert.NotNil(t, metricsService.ResponseSize)
	assert.NotNil(t, metricsService.ActiveConnections)

	// Verify that metrics are registered
	assert.NoError(t, testutil.CollectAndCompare(metricsService.RequestCount, nil))
	assert.NoError(t, testutil.CollectAndCompare(metricsService.RequestDuration, nil))
}
