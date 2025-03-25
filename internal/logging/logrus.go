package logging

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type LogEntry struct {
	Level   logrus.Level
	Message string
	Fields  logrus.Fields
}

// LogService handles structured logging
type LogService struct {
	logger  *logrus.Logger
	mu      sync.Mutex
	logChan chan LogEntry
}

var _ Logger = (*LogService)(nil)

// NewLogService initializes a new logging service
func NewLogService(config LoggingConfig) *LogService {
	logger := logrus.New()

	// Set up log rotation - @TODO: move to config or env variables
	fileLogger := &lumberjack.Logger{
		Filename:   "logs/api-gateway.log",
		MaxSize:    10, // Max megabytes before log is rotated
		MaxBackups: 3,  // Max number of old log files to keep
		MaxAge:     28, // Max number of days to retain old log files
		Compress:   true,
	}

	// Write logs to both file and stdout
	logger.SetOutput(io.MultiWriter(fileLogger, os.Stdout))

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel // Default to Info level if parsing fails
	}
	logger.SetLevel(level)

	// Set log format
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	// Read buffered channel size from environment variable
	logChanSize := 1000 // Default value
	if envLogChanSize := os.Getenv("LOG_BUFFERED_CHANNEL_SIZE"); envLogChanSize != "" {
		if size, err := strconv.Atoi(envLogChanSize); err == nil && size > 0 {
			logChanSize = size
		}
	}
	// Read batch size from environment variable
	batchSize := 5 // Default value

	if envBatchSize := os.Getenv("LOG_BATCH_SIZE"); envBatchSize != "" {
		if size, err := strconv.Atoi(envBatchSize); err == nil && size > 0 {
			batchSize = size
		}
	}

	logService := &LogService{
		logger:  logger,
		logChan: make(chan LogEntry, logChanSize), // Buffered channel for log entries
	}

	// Start a goroutine to process log entries asynchronously
	go logService.processLogs(batchSize)

	return logService
}

func (l *LogService) processLogs(batchSize int) {
	if batchSize <= 0 {
		batchSize = 5 // Default to 1000 if not set or invalid
	}
	batch := make([]LogEntry, 0, batchSize)
	ticker := time.NewTicker(5 * time.Second) // @TODO move the ticker timing to config or env variable
	defer ticker.Stop()
	for {
		select {
		case entry, ok := <-l.logChan:
			if !ok {
				// Channel is closed, flush remaining logs and exit
				if len(batch) > 0 {
					l.flushLogs(batch)
				}
				return
			}
			batch = append(batch, entry)
			if len(batch) >= 5 {
				l.flushLogs(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				l.flushLogs(batch)
				batch = batch[:0]
			}
		}
	}
}

func (l *LogService) flushLogs(batch []LogEntry) {
	for _, entry := range batch {
		l.logger.WithFields(entry.Fields).Log(entry.Level, entry.Message)
	}
}

// LogRequest logs information about a request
func (l *LogService) LogRequest(r *http.Request, duration time.Duration, status int, responseSize int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	fields := logrus.Fields{
		"timestamp":     time.Now().Format(time.RFC3339),
		"method":        r.Method,
		"path":          r.URL.Path,
		"remote_addr":   r.RemoteAddr,
		"duration_ms":   duration.Milliseconds(),
		"status":        status,
		"user_agent":    r.UserAgent(),
		"request_id":    r.Header.Get("X-Request-ID"),
		"response_size": responseSize,
	}

	level := logrus.InfoLevel
	message := "Request processed"

	if status >= 500 {
		level = logrus.ErrorLevel
		message = "Server error"
	} else if status >= 400 {
		level = logrus.WarnLevel
		message = "Client error"
	}

	// Send the log entry to the channel
	select {
	case l.logChan <- LogEntry{Level: level, Message: message, Fields: fields}:
	default:
		// Drop the log entry if the channel is full to avoid blocking
		fmt.Println("Log channel is full, dropping log entry")
	}
}

// Info logs an informational message
func (l *LogService) Info(msg string) {
	l.logger.Info(msg)
}

// Infof logs a formatted informational message
func (l *LogService) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Fatal logs a fatal message and exits the application
func (l *LogService) Fatal(msg string) {
	l.logger.Fatal(msg)
}

// Fatalf logs a formatted fatal message and exits the application
func (l *LogService) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *LogService) Shutdown() {
	close(l.logChan) // Close the channel to stop the goroutine
}
