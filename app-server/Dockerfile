# -----------------------------
# Build stage
# -----------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app-server

# Install build tools (musl-dev for static build)
RUN apk add --no-cache git build-base

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build statically
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# -----------------------------
# Run stage (minimal)
# -----------------------------
FROM scratch

WORKDIR /app-server

# Copy compiled binary
COPY --from=builder /app-server/main .

# Expose port
EXPOSE 8080

# Run binary
CMD ["/app-server/main"]
