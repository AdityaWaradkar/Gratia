FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./cmd/server


FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/auth-service /app/auth-service
RUN chmod +x /app/auth-service

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/auth-service"]
