package metricDaoIMpl

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type iDaoExport interface {
	Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error
	Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error)
	InstallDevice(c *gin.Context, device *deviceModels.DeviceVO) error
	UnInStallDevice(c *gin.Context, deviceId int64) error
}
