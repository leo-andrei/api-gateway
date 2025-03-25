package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds the application configuration
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
	Routes []Route `yaml:"routes"`
}

// Route represents a route configuration
type Route struct {
	Path        string `yaml:"path"`
	TargetURL   string `yaml:"targetUrl"`
	Method      string `yaml:"method"`
	RequireAuth bool   `yaml:"requireAuth"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(filename string) (*Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
