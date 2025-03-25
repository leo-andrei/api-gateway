package logging

import (
	"net/http"
	"time"
)

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Fatal(msg string)
	Fatalf(format string, args ...interface{})
	LogRequest(r *http.Request, duration time.Duration, status int, responseSize int)
	Shutdown()
}
