# ==============================
# Build Stage
# ==============================

FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN go build -o user_service ./cmd/main.go

# ==============================
# Runtime Stage
# ==============================

FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/user_service .

# Expose application port
EXPOSE 8081

# Run application
CMD ["./user_service"]