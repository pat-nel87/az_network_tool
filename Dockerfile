# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o az-network-analyzer .

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates graphviz

# Copy binary from builder to /usr/local/bin for global access
COPY --from=builder /app/az-network-analyzer /usr/local/bin/az-network-analyzer

# Create output directory
RUN mkdir -p /output && chmod 777 /output

WORKDIR /output

# Set entrypoint with absolute path
ENTRYPOINT ["/usr/local/bin/az-network-analyzer"]
CMD ["--help"]
