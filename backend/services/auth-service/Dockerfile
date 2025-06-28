# Dockerfile for auth-service
FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o main ./cmd/server/main.go

CMD ["./main"]
