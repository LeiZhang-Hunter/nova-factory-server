package devicemonitorservice

import (
	"context"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	deviceMonitorModel2 "nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
)

type DeviceMonitorService interface {
	List(c *gin.Context, req *devicemodels.DeviceListReq) (*devicemodels.DeviceInfoListData, error)
	Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error)
	Predict(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error)
	PredictQuery(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error)
	// DeviceLayout 设备布局
	DeviceLayout(c *gin.Context, floorId int64) (*deviceMonitorModel2.DeviceLayoutData, error)
	// ControlStatus 查询控制下发状态
	ControlStatus(c context.Context, req *deviceMonitorModel2.ControlStatusReq) (*deviceMonitorModel2.ControlStatusRes, error)
	// Control 设备控制
	Control(c *gin.Context, req *deviceMonitorModel2.ControlReq) (*deviceMonitorModel2.ControlRes, error)
	// GetRealTimeInfo 获取设备实时信息
	GetRealTimeInfo(c *gin.Context, deviceId uint64) (*devicemodels.DeviceVO, error)
	// GetRealTimeInfoList 批量获取设备实时信息
	GetRealTimeInfoList(c *gin.Context, deviceIds []string) ([]*devicemodels.DeviceVO, error)
}
