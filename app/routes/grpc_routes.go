package routes

import (
	"nova-factory-server/app/business/iot/metric/device/metriccontroller"

	"github.com/google/wire"
	"google.golang.org/grpc"
)

var GrpcProviderSet = wire.NewSet(NewGrpcEngine)

func NewGrpcEngine(server *metriccontroller.MetricServer) *grpc.Server {
	// 创建grpc服务
	s := grpc.NewServer()
	server.Metric.PrivateGrpcRoutes(s)
	server.CameraGrpc.PrivateGrpcRoutes(s)
	return s
}
