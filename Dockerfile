# Multi-stage Dockerfile for Vesper
# - Builds a statically linked Go binary in a builder image
# - Produces a small Alpine-based runtime image
# - Exposes port 8080 and mounts /data as a persistent volume for SQLite DB
#
# Notes:
# - The repository expects the SQLite DB file at ./data/tasks.db
#   In production, mount a host volume at /data so the DB persists outside the container.

# Builder stage

FROM golang:1.24 AS builder

# Working directory inside the builder
WORKDIR /src

# Use Go modules
COPY go.mod go.sum ./
RUN go env -w GOFLAGS=-mod=mod && go mod download

# Copy the rest of the source
COPY . .

# Build a small, optimized static binary.
# CGO_ENABLED=0 produces a static binary where possible.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /vesper ./cmd/server

# Runtime stage
FROM alpine:3.18

# Add ca-certificates for outbound TLS (required for OAuth / Google API)
RUN apk add --no-cache ca-certificates tzdata

# Create unprivileged user
RUN addgroup -S app && adduser -S -G app app

# Copy binary from builder
COPY --from=builder /vesper /usr/local/bin/vesper

# Create data directory for SQLite DB and give ownership to the app user.
RUN mkdir -p /data \
    && chown app:app /data

# Use non-root user for running
USER app
WORKDIR /home/app

# Configuration via environment variables
ENV PORT=8080
ENV DATA_DIR=/data

# Expose the default port
EXPOSE 8080

# Declare volume for persistent DB storage
VOLUME [ "/data" ]

# Default command
CMD [ "/usr/local/bin/vesper" ]
