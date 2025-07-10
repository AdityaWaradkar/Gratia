# Build Stage
FROM golang:1.24.4 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary from main.go inside cmd/server/
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/server/main.go

# Final Stage (minimal image)
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from builder
COPY --from=builder /app/auth-service .

# Expose port 8081 (your service port)
EXPOSE 8081

# Run the binary
CMD ["./auth-service"]
