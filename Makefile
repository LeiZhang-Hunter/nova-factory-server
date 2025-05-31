

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


