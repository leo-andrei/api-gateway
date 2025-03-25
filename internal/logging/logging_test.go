package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogService(t *testing.T) {
	cfg := LoggingConfig{
		Level:  "info",
		Format: "json",
	}
	logService := NewLogService(cfg)

	assert.NotNil(t, logService)
	assert.Equal(t, "info", logService.logger.GetLevel().String())
}
