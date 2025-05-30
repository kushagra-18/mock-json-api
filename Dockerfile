# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 for a static build, useful for alpine images
# -ldflags="-w -s" to strip debug information and reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/main ./cmd/server/main.go

# Stage 2: Create the runtime image
FROM alpine:latest

WORKDIR /app

# Copy the compiled application binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the configuration file
# Ensure config.yaml is suitable for default deployment or manage via env vars
COPY config.yaml /app/config.yaml

# Expose the port the application runs on (ensure this matches config)
# The port should ideally be fetched from the config file or an env var at runtime if it can change.
# For Dockerfile, EXPOSE is more of a documentation / hint.
# The actual port binding happens during `docker run -p <host_port>:<container_port>`.
# Assuming default port 8080 from config.yaml.
EXPOSE 8080

# Command to run the application
CMD ["/app/main"]
