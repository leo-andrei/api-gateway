# API Gateway with Metrics Logging

A lightweight API gateway that relays requests while providing comprehensive metrics logging and monitoring. This microservice focuses on capturing detailed performance data for all relayed requests.

## Features

- **Lightweight API Gateway**: Simple request relay with minimal routing logic
- **Comprehensive Metrics**: Captures timestamps, response times, endpoint usage, and request counts
- **Prometheus Integration**: Exposes metrics in Prometheus format for easy monitoring
- **Structured Logging**: JSON-formatted logs for easy parsing and analysis
- **Authentication Support**: Optional JWT authentication for protected routes
- **Containerized**: Ready to deploy with Docker and Docker Compose
- **Configurable**: External YAML configuration for routes and settings
- **Modular Design**: Clean separation of concerns for maintainability
- **Graceful Shutdown**: Handles termination signals and ensures proper cleanup
- **Extensible Logging and Metrics**: Interfaces for logging and metrics allow easy integration with other providers.

## Architecture

The gateway follows a modular architecture with clear separation of concerns:

```
api-gateway/
├── internal/             # Internal packages
│   ├── gateway/          # Core gateway functionality
│   ├── metrics/          # Metrics collection
│   │   ├── metrics.go    # Metrics interface definition
│   │   ├── prometheus.go # Prometheus-based implementation of the Metrics interface
│   ├── logging/          # Logging services
│   │   ├── logger.go     # Logger interface definition
│   │   ├── logrus.go     # Logrus-based implementation of the Logger interface
│   └── middleware/       # HTTP middleware
├── pkg/                  # Public packages
├── config.yaml           # Configuration file
├── Dockerfile            # Dockerfile for containerization
├── docker-compose.yml    # Docker Compose configuration
├── main.go               # Entry point to the application
└── README.md             # Documentation
```

## Metrics Captured

The gateway captures the following metrics:

- **Request Counts**: Total number of requests by path, method, and status code
- **Request Duration**: Response time in seconds (histogram)
- **Request Size**: Size of incoming requests in bytes
- **Response Size**: Size of outgoing responses in bytes
- **Active Connections**: Number of currently active connections

## Logging

The logging system is implemented using the `Logger` interface, which allows for easy integration with different logging providers. The default implementation uses `logrus`.

### Key Features:
- **Log Rotation**: Logs are rotated using the `lumberjack` package. Logs expire after a configurable number of days and are compressed.
- **Buffered Logging**: Logs are processed asynchronously using a buffered channel (default size: 1000).
- **Batch Processing**: Logs are written in batches (default size: 5) to reduce the number of writes.

### Extensibility:
- The `Logger` interface is defined in `internal/logging/logger.go`.
- The default implementation (`LogService`) uses `logrus` and is located in `internal/logging/logrus.go`.
- You can add new logging implementations (e.g., `ZapLogger`) by implementing the `Logger` interface.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose

### Running Locally

1. Clone the repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Run the gateway:
   ```
   go run main.go
   ```

### Running with Docker

Build and run the Docker container:

```
docker build -t api-gateway .
docker run -p 8080:8080 api-gateway
```

### Running with Docker Compose

1. Ensure `docker-compose.yml` is configured correctly.
2. Start the services:
   ```
   docker-compose up --build -d
   ```
   
## Test requests

For testing together with a service, I created a simple one at https://github.com/leo-andrei/user-service 
You can clone this repo in the same parent folder with this service, make the build with docker compose and test it with: 
```
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" http://localhost:8080/api/users
```

## Configuration

The gateway is configured using a YAML file (`config.yaml`) located in the root of the project. Example configuration:

```yaml
server:
  port: 8080

logging:
  level: info
  format: json

routes:
  - path: "/api/users"
    targetUrl: "http://user-service:8081/users"
    method: "GET"
    requireAuth: true
  
  - path: "/api/products"
    targetUrl: "http://product-service:8082/products"
    method: "GET"
    requireAuth: false
```

### Environment Variables for Logging

The logging system supports the following environment variables for configuration:
- `LOG_BUFFERED_CHANNEL_SIZE`: Size of the buffered channel for log entries (default: 1000).
- `LOG_BATCH_SIZE`: Number of log entries to process in a single batch (default: 5).

Example:
```bash
export LOG_BUFFERED_CHANNEL_SIZE=2000
export LOG_BATCH_SIZE=10
```

## Monitoring

The gateway exposes metrics in Prometheus format at the `/metrics` endpoint. You can use Prometheus to scrape these metrics.
To access the Prometheus UI use: http://localhost:9090/ and run your queries.

### Example Prometheus Configuration

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'api-gateway'
    static_configs:
      - targets: ['host.docker.internal:8080'] # Use host.docker.internal for Docker
```

## Testing

Run the tests:

```
go test ./...
```

## Adding Authentication

The gateway includes a simple JWT authentication middleware. To enable authentication for a route, set `requireAuth: true` in the route configuration.

In a production environment, you would extend the `AuthMiddleware` to validate JWT tokens properly.

## Extending the Gateway

### Adding a New Route

Add a new route to the `config.yaml` file:

```yaml
routes:
  - path: "/api/new-endpoint"
    targetUrl: "http://new-service:8084/resource"
    method: "POST"
    requireAuth: true
```

### Adding Custom Middleware

Create a new middleware in the `internal/middleware` directory and add it to the middleware chain in `gateway.go`.

## Graceful Shutdown

The gateway handles termination signals (`SIGINT` and `SIGTERM`) to ensure a graceful shutdown. This includes:
- Closing active connections.
- Flushing logs and metrics.
- Cleaning up resources.

## Known Limitations

- Static service discovery: The gateway currently uses static routes defined in `config.yaml`. Dynamic service discovery (e.g., via Consul or Kubernetes) can be added for scalability.
- Limited authentication: The JWT middleware is basic and should be extended for production use.
- No distributed tracing: Consider integrating tools like Jaeger or Zipkin for tracing requests across services.

## Future Improvements

- Add support for dynamic service discovery.
- Integrate distributed tracing for better observability.
- Enhance logging with centralized log aggregation (e.g., ELK Stack or Loki).
- Add circuit breakers and retries for fault tolerance.

## Good to Know About Logs

- For testing purposes, the logs are currently displayed both in a file and stdout.
- For high throughput, the following solutions are implemented:
  1. **Log Rotation**: Using the `lumberjack` package, logs expire after a configurable number of days and are compressed.
  2. **Buffered Channel**: A buffered channel (default size: 1000) is used to keep logs and process them asynchronously with a goroutine.
  3. **Batch Processing**: Logs are processed in batches (default size: 5) to limit the number of writes.


