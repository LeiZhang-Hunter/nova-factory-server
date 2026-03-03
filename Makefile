

init:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/swaggo/swag/cmd/swag@latest
wire:
	cd app/ && wire

swag:
	cd app/ && swag  init

# Parsing protobuf files and generating go files.
pb:
	protoc --go_out=. --go-grpc_out=. ./manifest/protobuf/metric/server.proto
	# 生成message
	#protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative ./manifest/protobuf/metric/server.proto
	# 生成grpc service
	#protoc --proto_path=proto --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/echo.proto
	#protoc -I. --go_out=/app/pkg/metric/grpc/v1 --go-grpc_out=. manifest/protobuf/metric/server.proto


# Docker Registry variables
DOCKER_REGISTRY := crpi-wvgmi16z1eaqbhcu.cn-shanghai.personal.cr.aliyuncs.com
DOCKER_NAMESPACE := novawatcherio
IMAGE_NAME := nova-factory-server
IMAGE_TAG := latest
FULL_IMAGE_NAME := $(DOCKER_REGISTRY)/$(DOCKER_NAMESPACE)/$(IMAGE_NAME):$(IMAGE_TAG)

# Build Docker image
docker-build:
	@echo "==> Make Build completed!"
	docker build -t $(FULL_IMAGE_NAME) .

# Push Docker image to registry
docker-push:
	@echo "==> Make Push completed!"
	docker push $(FULL_IMAGE_NAME)

# Build and push all at once
docker-build-push: docker-build docker-push
	@echo "==> Make Build & Push completed!"
