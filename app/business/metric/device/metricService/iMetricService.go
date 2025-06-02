package metricService

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
)

type IMetricService interface {
	Export(c context.Context, request *v1.ExportMetricsServiceRequest) error
}
