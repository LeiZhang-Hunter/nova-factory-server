package metricDao

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type IMetricDao interface {
	Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error
	Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error)
}
