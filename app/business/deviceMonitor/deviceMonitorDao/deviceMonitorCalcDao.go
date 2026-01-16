package deviceMonitorDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"
)

// DeviceMonitorCalcDao 统计监控设备
type DeviceMonitorCalcDao interface {
	// CounterByTimeRange 通过时间范围counter 设备写入总数
	CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricModels.MetricQueryData, error)
	// CounterByDevice 通过时间范围counter 根据设备分组记录设备写入总数
	CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*deviceMonitorModel.TypeDeviceCounterRank, error)
}
