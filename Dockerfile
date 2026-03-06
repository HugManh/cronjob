# ─── Stage 1: Builder ─────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

# Install git for fetching modules that require VCS
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies first for faster rebuilds
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/cronjob \
    ./cmd/api

# ─── Stage 2: Runner ──────────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/cronjob .

# Copy static assets needed at runtime
COPY --from=builder /app/public ./public
COPY --from=builder /app/views  ./views

# Use the non-root user
USER appuser

EXPOSE 8080

ENTRYPOINT ["./cronjob"]
