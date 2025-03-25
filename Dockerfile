# filepath: /Users/leonteandrei/Documents/testProjects/api-gateway/docker/Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o api-gateway ./main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/api-gateway .

# Copy the config.yaml file
COPY config.yaml ./config.yaml

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["./api-gateway"]