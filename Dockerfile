# Multi-stage Dockerfile for Paperless MCP Server
# Stage 1: Build the Go binary
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Set the environment variables for Go cross-compilation
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build \
    -ldflags="-w -s" \
    -o paperless-mcp \
    ./cmd/server

# Stage 2: Create minimal runtime image
FROM --platform=$BUILDPLATFORM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 paperless && \
    adduser -D -u 1000 -G paperless paperless

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/paperless-mcp .

# Change ownership
RUN chown -R paperless:paperless /app

# Switch to non-root user
USER paperless

# Expose HTTP port
EXPOSE 8080

# Health check (for HTTP mode)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the server
ENTRYPOINT ["/app/paperless-mcp"]
