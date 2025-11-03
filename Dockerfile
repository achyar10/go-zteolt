# Multi-stage build untuk Go application
# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Tools yang dibutuhkan saat build
RUN apk add --no-cache git ca-certificates tzdata

# Copy & download dependencies lebih awal agar cache build efektif
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Copy seluruh sumber
COPY . .

# Pastikan folder templates ada (agar COPY di stage runtime tidak gagal)
RUN mkdir -p /app/templates

# Build arguments untuk version dan build time
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
# Untuk buildx (multi-arch)
ARG TARGETOS
ARG TARGETARCH

# Build server binary (static, kecil, dan reproducible-ish)
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -buildvcs=false \
    -ldflags='-w -s -extldflags "-static" -X main.version='${VERSION}' -X main.buildTime='${BUILD_TIME}' -X main.gitCommit='${GIT_COMMIT} \
    -o bin/server cmd/server/main.go

# Stage 2: Runtime stage
FROM alpine:3.20

# Install runtime deps (certificates, tzdata, curl untuk healthcheck)
RUN apk --no-cache add ca-certificates tzdata curl

# Non-root user demi keamanan
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary & assets dari builder
COPY --from=builder /app/bin/server ./bin/server
COPY --from=builder /app/templates ./templates

# Ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Port default aplikasi
EXPOSE 8080

# Env
ENV GIN_MODE=release
# Opsional: set zona waktu
# ENV TZ=Asia/Jakarta

# Build args perlu dideklarasikan ulang di stage final agar LABEL dapat nilainya
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Labels
LABEL maintainer="Achyar Anshorie <achyar@matik.id>" \
    org.opencontainers.image.title="ZTE OLT Management API" \
    org.opencontainers.image.description="ZTE OLT Management API" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.revision="${GIT_COMMIT}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.source="https://github.com/achyar10/go-zteolt"

# Jalankan server
ENTRYPOINT ["./bin/server"]