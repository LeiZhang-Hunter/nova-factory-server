package devicemonitordao

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
)

// DeviceMonitorCalcDao 统计监控设备
type DeviceMonitorCalcDao interface {
	// CounterByTimeRange 通过时间范围counter 设备写入总数
	CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricmodels.MetricQueryData, error)
	// CounterByDevice 通过时间范围counter 根据设备分组记录设备写入总数
	CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error)
}
