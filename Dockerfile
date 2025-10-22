# Multi-stage build untuk Go application
# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git dan package yang dibutuhkan
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build arguments untuk version dan build time
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build server binary dengan optimization
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -X main.version=${VERSION} \
    -X main.buildTime=${BUILD_TIME} \
    -X main.gitCommit=${GIT_COMMIT} \
    -a -installsuffix cgo \
    -o bin/server cmd/server/main.go

# Build CLI binary (optional)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -X main.version=${VERSION} \
    -X main.buildTime=${BUILD_TIME} \
    -X main.gitCommit=${GIT_COMMIT} \
    -a -installsuffix cgo \
    -o bin/cli main.go

# Stage 2: Runtime stage
FROM alpine:latest

# Install ca-certificates untuk HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user untuk security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary dari builder stage
COPY --from=builder /app/bin/server ./bin/server
COPY --from=builder /app/bin/cli ./bin/cli

# Copy templates directory (jika ada)
COPY --from=builder /app/templates ./templates

# Change ownership ke non-root user
RUN chown -R appuser:appgroup /app

# Switch ke non-root user
USER appuser

# Expose port (default dari aplikasi)
EXPOSE 8080

# Environment variables
ENV GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Default command untuk menjalankan server
CMD ["./bin/server"]

# Labels untuk metadata
LABEL maintainer="Achyar Anshorie <achyar@matik.id>" \
      version="${VERSION}" \
      description="ZTE OLT Management API" \
      org.opencontainers.image.source="https://github.com/achyar10/go-zteolt"