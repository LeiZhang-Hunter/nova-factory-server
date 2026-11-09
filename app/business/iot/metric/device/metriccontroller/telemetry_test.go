package metriccontroller

import (
	"context"
	"errors"
	"testing"

	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	metricservice "nova-factory-server/app/business/iot/metric/device/metricservice"
)

type failingTelemetryService struct {
	metricservice.IMetricService
	calls int
}

func (s *failingTelemetryService) Export(context.Context, *v1.ExportMetricsServiceRequest) error {
	s.calls++
	return errors.New("storage unavailable")
}

func TestTelemetryFailsClosedAndPropagatesStorageFailure(t *testing.T) {
	service := &failingTelemetryService{}
	controller := NewMetric(service, nil)
	request := &v1.ExportMetricsServiceRequest{}
	if _, err := controller.ReportDeviceInfo(context.Background(), request); err == nil || service.calls != 0 {
		t.Fatal("missing authorizer allowed write")
	}
	controller.ValidateGateway = func(context.Context, *v1.ExportMetricsServiceRequest) error {
		return errors.New("unauthorized")
	}
	if _, err := controller.ReportDeviceInfo(context.Background(), request); err == nil || service.calls != 0 {
		t.Fatal("unauthorized write")
	}
	controller.ValidateGateway = func(context.Context, *v1.ExportMetricsServiceRequest) error {
		return nil
	}
	if reply, err := controller.ReportDeviceInfo(context.Background(), request); err == nil || reply.Code == 0 || service.calls != 1 {
		t.Fatal("storage failure acknowledged")
	}
	endpoint := &gatewayTelemetryServer{
		metric: controller,
	}
	if _, err := endpoint.ReportTimeData(context.Background(), &v1.ExportTimeDataRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatal("legacy write RPC exposed")
	}
	if _, err := endpoint.ReportScheduleLog(context.Background(), &v1.ExportControlLogRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatal("legacy log RPC exposed")
	}
}
