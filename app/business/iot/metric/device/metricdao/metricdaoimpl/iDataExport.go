package metricdaoimpl

import (
	"context"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
)

type iDaoExport interface {
	Export(ctx context.Context, data []*metricmodels.NovaMetricsDevice) error
	ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error
	Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error)
	InstallDevice(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData) error
	UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error
	// InstallRunStatusDevice 运行状态设备模板
	InstallRunStatusDevice(c *gin.Context, deviceId int64) error
	// UnInStallRunStatusDevice 卸载设备运行状态模板
	UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error
	Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error)
	List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error)
	Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error)
	Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error)
}
