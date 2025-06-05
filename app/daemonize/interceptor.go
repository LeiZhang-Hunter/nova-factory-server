package daemonize

import (
	"context"
	"errors"
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
	val, ok := md["company_uuid"]
	if !ok {
		zap.L().Error("get grpc request meta data company_uuid not exist")
		return nil, errors.New("grpc request CodeNotAuthorized")
	}

	zap.L().Info("get company_uuid:%s", zap.Any("val", val))
	//companyUuid := val[0]
	//ctx = context.WithValue(ctx, common.Cid, companyUuid)
	return handler(ctx, req)
}
