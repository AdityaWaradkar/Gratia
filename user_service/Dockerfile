# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install certificates
RUN apk add --no-cache ca-certificates

# Copy dependency files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o user_service ./cmd

# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /app

# Install certificates
RUN apk add --no-cache ca-certificates wget

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary
COPY --from=builder /app/user_service .

# Use non-root user
USER appuser

# Expose application port
EXPOSE 8081

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider http://localhost:8081/health || exit 1

# Run application
CMD ["./user_service"]