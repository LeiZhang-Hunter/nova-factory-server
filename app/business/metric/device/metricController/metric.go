package metricController

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	v1 "nova-factory-server/app/pkg/metric/grpc/v1"
)

type Metric struct {
	v1.UnimplementedDeviceReportServiceServer
}

func NewMetric() *Metric {
	return &Metric{}
}

func (m *Metric) PrivateRoutes(router *grpc.Server) {
	v1.RegisterDeviceReportServiceServer(router, m)
}

func (m *Metric) ReportContainer(context.Context, *v1.DeviceData) (*v1.NodeRes, error) {

	fmt.Println("article 进来了")
	return &v1.NodeRes{
		Code: 0,
	}, nil
}
