package metricController

import (
	"context"
	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"google.golang.org/grpc"
	"nova-factory-server/app/business/metric/device/metricService"
)

type Metric struct {
	service metricService.IMetricService
	v1.UnimplementedDeviceReportServiceServer
}

func NewMetric(service metricService.IMetricService) *Metric {
	return &Metric{
		service: service,
	}
}

func (m *Metric) PrivateGrpcRoutes(router *grpc.Server) {
	v1.RegisterDeviceReportServiceServer(router, m)
}

func (m *Metric) PrivateRoutes(router *gin.RouterGroup) {
}

func (m *Metric) ReportDeviceInfo(c context.Context, request *v1.ExportMetricsServiceRequest) (*v1.NodeRes, error) {
	err := m.service.Export(c, request)
	if err != nil {
		return &v1.NodeRes{
			Code: -1,
		}, err
	}
	return &v1.NodeRes{
		Code: 0,
	}, nil
}
