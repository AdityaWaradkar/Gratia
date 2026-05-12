# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth_service ./cmd

# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth_service .

COPY .env .

EXPOSE 8080

CMD ["./auth_service"]