package metricdao

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
)

// IMetricAdderDao 添加器
type IMetricAdderDao interface {
	// ExportTimeData 导出时序数据
	ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error
	Export(ctx context.Context, data []*entity.NovaMetricsDevice) error
}
