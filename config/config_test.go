package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_ValidFile(t *testing.T) {
	configPath := "testdata/valid_config.yaml"
	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	configPath := "testdata/invalid_config.yaml"
	_, err := LoadConfig(configPath)
	assert.Error(t, err)
}

func TestLoadConfig_MissingFile(t *testing.T) {
	configPath := "testdata/missing_config.yaml"
	_, err := LoadConfig(configPath)
	assert.Error(t, err)
}
