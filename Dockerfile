# ==========================================
# STAGE 1: Build Binary
# ==========================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary yang dioptimasi
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/main.go

# ==========================================
# STAGE 2: Minimal Runtime Environment
# ==========================================
FROM alpine:3.20

WORKDIR /app

# Install CA Certificates (Wajib untuk HTTPS Google OAuth & TLS SMTP)
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta

# Buat folder uploads, configs, dan migrations
RUN mkdir -p /app/uploads /app/internal/configs /app/scripts/migrations

# Copy binary dan file migrasi dari builder
COPY --from=builder /app/server .
COPY --from=builder /app/scripts/migrations ./scripts/migrations

# Buat user non-root untuk keamanan kontainer
RUN adduser -D -u 1001 appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 9888

CMD ["./server"]
