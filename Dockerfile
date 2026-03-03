# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory directly inside the container
WORKDIR /app

# Enable CGO off for a static build, and set proxies if necessary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct

# Copy go.mod and go.sum first to take advantage of Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application from the app directory where main.go is located
RUN go build -o baize ./app/

# Stage 2: Create a minimal runtime image
FROM alpine:latest
# Install CA certificates for external HTTPS requests or timezone data if needed
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /build

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/baize /build/baize

# Copy necessary configuration, template and other static files
COPY --from=builder /app/template /build/template
COPY --from=builder /app/config /build/config

# Expose port (80 or 8080 according to config, we use 8080 here based on config.yaml, 
# although old dockerfile said 80, you can map it later)
EXPOSE 8080 10050

ENTRYPOINT ["./baize"]
