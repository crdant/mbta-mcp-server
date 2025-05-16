FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=$(git describe --tags --always)" -o mbta-mcp-server ./cmd/server

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Import ca-certificates for secure connections
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from build stage
COPY --from=builder /app/mbta-mcp-server .

# Expose port
EXPOSE 8080

# Run the service
CMD ["./mbta-mcp-server"]