package deviceMonitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type DeviceMonitorService interface {
	List(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error)
	Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error)
}
