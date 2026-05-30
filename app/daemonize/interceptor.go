package daemonize

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"nova-factory-server/app/constant/agent"
)

func CompanyValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)

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
		gateway_id, ok = md["gatewayid"]
		if !ok {
			return nil, errors.New("grpc request CodeNotAuthorized")
		}
	}

	ctx = context.WithValue(ctx, agent.USERNAME, username)
	ctx = context.WithValue(ctx, agent.PASSWORD, pasword)
	ctx = context.WithValue(ctx, agent.GATEWAYID, gateway_id)

	return handler(ctx, req)
}
