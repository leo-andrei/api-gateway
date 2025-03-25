package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/leo-andrei/api-gateway/config"
	"github.com/leo-andrei/api-gateway/internal/logging"
	"github.com/leo-andrei/api-gateway/internal/metrics"
	"github.com/leo-andrei/api-gateway/internal/middleware"
)

// Gateway represents the API gateway
type Gateway struct {
	config         *config.Config
	router         *mux.Router
	server         *http.Server
	logService     logging.Logger
	metricsService metrics.Metrics
}

// NewGateway initializes a new API gateway
func NewGateway(cfg *config.Config, logger logging.Logger, metrics metrics.Metrics) *Gateway {
	router := mux.NewRouter()

	return &Gateway{
		config:         cfg,
		router:         router,
		logService:     logger,
		metricsService: metrics,
	}
}

// SetupRoutes configures the routes for the gateway
func (g *Gateway) SetupRoutes() {
	// Add metrics endpoint
	g.router.Handle("/metrics", promhttp.Handler())

	// Add health check endpoint
	g.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}).Methods("GET")

	// Configure routes from config
	for _, route := range g.config.Routes {
		// Create handler with auth middleware if required
		handler := http.Handler(CreateProxyHandler(route))
		if route.RequireAuth {
			handler = middleware.AuthMiddleware(handler)
		}

		// Apply metrics middleware
		handler = middleware.MetricsMiddleware(handler, g.metricsService, g.logService)

		// Register route
		g.router.Handle(route.Path, handler).Methods(route.Method)
	}
}

// Run starts the gateway server
func (g *Gateway) Run() error {
	g.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", g.config.Server.Port),
		Handler: g.router,
	}

	return g.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (g *Gateway) Shutdown(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}
