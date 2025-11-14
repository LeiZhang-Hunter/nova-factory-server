package metricDao

import (
	"context"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type IControlLogDao interface {
	Export(ctx context.Context, data []*metricModels.NovaControlLog) error
}
