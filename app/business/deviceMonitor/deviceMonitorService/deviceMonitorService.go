package deviceMonitorService

import (
	"context"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"

	"github.com/gin-gonic/gin"
)

type DeviceMonitorService interface {
	List(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error)
	Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error)
	Predict(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error)
	PredictQuery(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error)
	// DeviceLayout 设备布局
	DeviceLayout(c *gin.Context, floorId int64) (*deviceMonitorModel.DeviceLayoutData, error)
	// ControlStatus 查询控制下发状态
	ControlStatus(c context.Context, req *deviceMonitorModel.ControlStatusReq) (*deviceMonitorModel.ControlStatusRes, error)
}
