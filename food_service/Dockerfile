# Build stage for compiling the food service binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install CA certificates and timezone data for secure connections and accurate time handling
RUN apk add --no-cache ca-certificates tzdata

# Copy dependency files first to leverage Docker caching
COPY go.mod go.sum ./

# Download all Go module dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the food service binary with optimizations (statically linked, stripped)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o food_service ./cmd/server

# Runtime stage for running the food service
FROM alpine:latest

WORKDIR /app

# Import necessary certificates and timezone data directly from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built binary from the builder stage
COPY --from=builder /app/food_service .

# Switch to the non-root user
USER appuser

# Expose the default service port
EXPOSE 8082

# Health check endpoint to monitor service status
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8082/health || exit 1

# Start the food service
CMD ["./food_service"]