ARG GO_VERSION=1.26.1
ARG BUILD_TAGS

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src

# Copy dependency files first for layer caching
COPY nova-factory-server/go.mod nova-factory-server/go.sum ./

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
RUN go mod download

# Copy the server project
COPY nova-factory-server/ .

# Copy symlinked add-on modules (symlinks can't cross Docker build context boundary)
COPY nova-factory-addons-be/shop app/business/shop
COPY nova-factory-addons-be/datasyncapi app/business/datasyncapi

# Build the server binary from app/ directory
# Build tags: ai iot shop datasyncapi (matches production deployment)
WORKDIR /src/app
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -tags="${BUILD_TAGS}" \
    -o /out/nova-factory-server .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        tzdata \
    && rm -rf /var/lib/apt/lists/*

# Set timezone
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# Copy binary
COPY --from=builder /out/nova-factory-server /app/nova-factory-server

# Create directories for upload files and logs
RUN mkdir -p /app/public /app/private /app/logs

EXPOSE 8080 10050

ENTRYPOINT ["/app/nova-factory-server"]
CMD ["--config=/app/config/config.yaml"]
