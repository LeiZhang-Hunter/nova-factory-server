package metriccontroller

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
)

// gatewayTelemetryServer exposes only authorized numeric telemetry.
// Legacy time-data and log RPCs remain unimplemented on this endpoint.
type gatewayTelemetryServer struct {
	v1.UnimplementedDeviceReportServiceServer
	metric *Metric
}

func newGatewayTelemetryServer(metric *Metric) *gatewayTelemetryServer {
	return &gatewayTelemetryServer{
		metric: metric,
	}
}

func (s *gatewayTelemetryServer) ReportDeviceInfo(ctx context.Context, request *v1.ExportMetricsServiceRequest) (*v1.NodeRes, error) {
	return s.metric.ReportDeviceInfo(ctx, request)
}
