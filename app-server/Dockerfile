# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app-server

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main

# Run stage
FROM alpine:latest

WORKDIR /app-server

COPY --from=builder /app-server/main .
COPY --from=builder /app-server/.env .env

EXPOSE 8080

CMD ["./main"]
