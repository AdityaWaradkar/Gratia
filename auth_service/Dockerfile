# -------- Stage 1: Build the Go binary --------
FROM golang:1.22-alpine AS builder

# Install git (needed for go mod download in alpine)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the whole project
COPY . .

# Build the binary (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./cmd

# -------- Stage 2: Create a minimal runtime image --------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/auth-service .

# Copy env file (optional — recommended to mount instead)
# COPY .env .env

EXPOSE 8080

# Run the binary
USER nonroot:nonroot
ENTRYPOINT ["./auth-service"]
