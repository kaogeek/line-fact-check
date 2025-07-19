# Multi-stage build for Go workspace executables

# Build stage
FROM golang:1.24.4-alpine AS builder

# Install git (needed for Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /workspace

RUN mkdir bin

# Copy Go workspace files
COPY go.work go.work.sum ./

# Copy all modules
COPY . ./

# Download dependencies for the workspace
RUN go mod download

# Build factcheck API
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/factcheck-api ./factcheck/cmd/api

# Build foo API  
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/foo-api ./foo/cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy built executables from builder stage
COPY --from=builder /workspace/bin/* .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports (both services use 8080 by default)
EXPOSE 8080

# Default command runs factcheck-api
CMD ["./factcheck-api"]
