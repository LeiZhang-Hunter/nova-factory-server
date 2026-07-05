

init:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Addon usage:
#   make wire                         # regenerate Wire without linking external addons
#   make wire-addons                  # link all addons from ADDONS_DIR, then regenerate Wire
#   make wire-addons ADDONS="shop erp" # link only selected addons, then regenerate Wire
#   make wire-addons ADDONS_DIR=/path/to/nova-factory-addons-be
#
# ADDONS_DIR points to the external addon repository.
# ADDONS defaults to all non-hidden first-level directories under ADDONS_DIR.
# WIRE_TAGS are always enabled for the main server; ADDONS are appended by wire-addons.
ADDONS_DIR ?= /home/zhanglei/project/nova/nova-factory-addons-be
ADDONS ?= $(shell if [ -d "$(ADDONS_DIR)" ]; then find "$(ADDONS_DIR)" -mindepth 1 -maxdepth 1 -type d ! -name ".*" -exec basename {} \; | sort | tr '\n' ' '; fi)
WIRE_TAGS ?= ai iot

addons:
	go run ./tools/addonsync -addons-dir "$(ADDONS_DIR)" -enabled "$(ADDONS)" -link

wire:
	go run ./tools/addonsync
	cd app/ && wire gen -tags "$(WIRE_TAGS)"

wire-addons: addons
	cd app/ && wire gen -tags "$(WIRE_TAGS) $(ADDONS)"

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


