# Build stage for compiling the claim service binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install CA certificates for secure connections
RUN apk add --no-cache ca-certificates

# Copy dependency files first to leverage Docker caching
COPY go.mod go.sum ./

# Download all Go module dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the claim service binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o claim_service ./cmd

# Runtime stage for running the claim service
FROM alpine:latest

WORKDIR /app

# Install CA certificates and wget for health checks
RUN apk add --no-cache ca-certificates wget

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built binary from the builder stage
COPY --from=builder /app/claim_service .

# Switch to the non-root user
USER appuser

# Expose the default service port
EXPOSE 8083

# Health check endpoint to monitor service status
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider http://localhost:8083/health || exit 1

# Start the claim service
CMD ["./claim_service"]