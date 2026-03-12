package metricController

import (
	"context"
	"nova-factory-server/app/business/iot/metric/device/metricService"

	v1 "github.com/novawatcher-io/nova-factory-payload/camera/v1"
	"google.golang.org/grpc"
)

type CameraGrpc struct {
	service metricService.ICameraService
	v1.UnimplementedCameraServiceServer
}

func NewCamera(service metricService.ICameraService) *CameraGrpc {
	return &CameraGrpc{
		service: service,
	}
}

func (c *CameraGrpc) PrivateGrpcRoutes(router *grpc.Server) {
	v1.RegisterCameraServiceServer(router, c)
}

func (c *CameraGrpc) Report(ctx context.Context, req *v1.CameraData) (*v1.CameraResponse, error) {
	err := c.service.Report(ctx, req)
	if err != nil {
		return &v1.CameraResponse{
			Code: -1,
		}, err
	}
	return &v1.CameraResponse{
		Code: 0,
	}, nil
}
