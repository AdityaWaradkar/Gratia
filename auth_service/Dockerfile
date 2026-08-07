# Build stage for compiling the highly optimized Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install CA certificates and tzdata for secure network calls and accurate timezone handling
RUN apk add --no-cache ca-certificates tzdata

# Copy dependency files first to maximize Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the core service code into the build container
COPY . .

# Compile a statically linked binary with debug symbols stripped for minimum size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o auth_service ./cmd/server

# Runtime stage using a pristine, minimal alpine image
FROM alpine:latest

WORKDIR /app

# Import necessary certificates and timezone data directly from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Establish a strict non-root user context to prevent privilege escalation vulnerabilities
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy only the finalized executable binary from the builder stage
COPY --from=builder /app/auth_service .

# Enforce the non-root user for all subsequent container operations
USER appuser

# Document the designated network port for container orchestration tools
EXPOSE 8080

# Actively monitor service health to allow orchestrators to restart frozen containers
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8080/health || exit 1

# Execute the isolated authentication service
CMD ["./auth_service"]