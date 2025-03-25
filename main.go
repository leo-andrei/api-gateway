package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leo-andrei/api-gateway/config"
	"github.com/leo-andrei/api-gateway/internal/gateway"
	"github.com/leo-andrei/api-gateway/internal/logging"
	"github.com/leo-andrei/api-gateway/internal/metrics"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml" // Default fallback
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger := logging.NewLogService(logging.LoggingConfig{Level: "error", Format: "text"})
		logger.Fatalf("Error loading config: %v", err)
	}

	// Initialize services
	logger := logging.NewLogService(cfg.Logging)
	metrics := metrics.NewMetricsService()

	// Create and run the gateway
	gw := gateway.NewGateway(cfg, logger, metrics)
	gw.SetupRoutes()

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Infof("Starting API Gateway on port %d", cfg.Server.Port)
		if err := gw.Run(); err != nil {
			logger.Fatalf("Error running gateway: %v", err)
		}
		logger.Infof("Started API Gateway on port %d", cfg.Server.Port)
	}()

	<-stop
	logger.Info("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gw.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down server: %v", err)
	}
	logger.Shutdown()

	logger.Info("API Gateway stopped")
}
