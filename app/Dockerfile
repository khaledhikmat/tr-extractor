# Stage 1: Build the Go binary
FROM golang:latest AS builder

# Set working directory
WORKDIR /app

# Cache dependencies first (optional, speeds up rebuilds)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Create a minimal image
FROM alpine:latest

# Install certificates (needed for many Go HTTP clients)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
