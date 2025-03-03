# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with version info
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api ./cmd/api


# Final stage
FROM alpine:3.18

ARG VERSION

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/api /app/api

# Set executable permissions
RUN chmod +x /app/api

# Create a non-root user to run the application
RUN adduser -D -u 1000 appuser
USER appuser

# Application metadata
LABEL maintainer="Your Name <your.email@example.com>"
LABEL version="${VERSION}"
LABEL description="E-commerce Go apps"

# Expose the application port
EXPOSE 8080

# Command to run
CMD ["/app/api"]