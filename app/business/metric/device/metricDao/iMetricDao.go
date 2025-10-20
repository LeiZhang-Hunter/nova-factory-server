package metricDao

import (
	"context"
	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type IMetricDao interface {
	Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error
	Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error)
	InstallDevice(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData) error
	UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error
	// InstallRunStatusDevice 运行状态设备模板
	InstallRunStatusDevice(c *gin.Context, deviceId int64) error
	// UnInStallRunStatusDevice 卸载设备运行状态模板
	UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error
	Predict(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error)
	List(c *gin.Context, req *deviceMonitorModel.DevDataReq) (*deviceMonitorModel.DevDataResp, error)
	Count(c *gin.Context, req *deviceMonitorModel.DevDataReq) (uint64, error)
	Query(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error)
	// ExportTimeData 导出时序数据
	ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error
}
