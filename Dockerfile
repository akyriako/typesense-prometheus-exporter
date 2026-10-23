# Stage 1: Build
FROM golang:1.26.10 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o typesense-prometheus-exporter ./cmd

# Stage 2: Package
FROM alpine:3.24.3

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/typesense-prometheus-exporter /app/typesense-prometheus-exporter

# Expose default metrics port
EXPOSE 8080

# Command to run the exporter
CMD ["/app/typesense-prometheus-exporter"]
