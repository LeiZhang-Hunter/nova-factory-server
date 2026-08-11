package metricdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
)

// IMetricQueryDao 查询器
type IMetricQueryDao interface {
	Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *entity.MetricQueryReq) (*entity.MetricQueryData, error)
	List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error)
	Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error)
	Query(c *gin.Context, req *entity.MetricDataQueryReq) (*entity.MetricQueryData, error)
	CounterByTimeRange(startTime int64, endTime int64, interval string) (*entity.MetricQueryData, error)
	CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error)
	StatDeviceStatus(c *gin.Context, startTime string, endTime string, status int) (*devicemonitormodel.DeviceStatusList, error)
	StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string, status int) (*devicemonitormodel.DeviceProcessList, error)
	StatDeviceRunStatus(c *gin.Context, startTime string, endTime string) ([]devicemonitormodel.DeviceRunStat, error)
	StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error)
	StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, interval string, status int) (*devicemonitormodel.DeviceProcessList, error)
	Metric(c *gin.Context, req *entity.MetricQueryReq) (*entity.MetricQueryData, error)
}
