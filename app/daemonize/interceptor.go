package daemonize

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CompanyValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	zap.L().Info("get md:%v", zap.Any("md", md))

	if !ok {
		zap.L().Error("get grpc request meta data failed")
		return nil, errors.New("grpc request CodeNotAuthorized")
	}
	username, ok := md["username"]
	if !ok {
		zap.L().Error("get grpc request meta data company_uuid not exist")
		return nil, errors.New("grpc request CodeNotAuthorized")
	}

	pasword, ok := md["password"]
	if !ok {
		return nil, errors.New("grpc request CodeNotAuthorized")
	}

	gateway_id, ok := md["gateway_id"]
	if !ok {
		return nil, errors.New("grpc request CodeNotAuthorized")
	}
	fmt.Println(username, pasword, gateway_id)

	//companyUuid := val[0]
	//ctx = context.WithValue(ctx, common.Cid, companyUuid)
	return handler(ctx, req)
}
