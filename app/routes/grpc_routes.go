package routes

import (
	"github.com/google/wire"
	"google.golang.org/grpc"
	"nova-factory-server/app/business/metric/device/metricController"
)

var GrpcProviderSet = wire.NewSet(NewGrpcEngine)

func NewGrpcEngine(server *metricController.MetricServer) *grpc.Server {
	// 创建grpc服务
	s := grpc.NewServer()
	server.Metric.PrivateGrpcRoutes(s)
	return s
}
