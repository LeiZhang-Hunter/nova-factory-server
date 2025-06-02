package metricDao

import (
	"context"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type IMetricDao interface {
	Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error
}
