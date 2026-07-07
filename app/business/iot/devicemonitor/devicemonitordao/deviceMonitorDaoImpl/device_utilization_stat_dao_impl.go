package deviceMonitorDaoImpl

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"

	"github.com/gin-gonic/gin"
)

// statRun 统计运行时设备
func (d *DeviceUtilizationDaoImpl) statDeviceStat(c *gin.Context, startTime string, endTime string,
	status int) (*devicemonitormodel.DeviceStatusList, error) {
	return d.metricDao.StatDeviceStatus(c, startTime, endTime, status)
}

// statDeviceProcess 统计设备运行过程
func (d *DeviceUtilizationDaoImpl) statDeviceProcess(c *gin.Context, startTime string, endTime string, interval string,
	status int) (*devicemonitormodel.DeviceProcessList, error) {
	return d.metricDao.StatDeviceProcess(c, startTime, endTime, interval, status)
}

// statDeviceProcess 统计设备运行过程
func (d *DeviceUtilizationDaoImpl) statDeviceRunStat(c *gin.Context, startTime string,
	endTime string) ([]devicemonitormodel.DeviceRunStat, error) {
	return d.metricDao.StatDeviceRunStatus(c, startTime, endTime)
}

// statDeviceStatByDeviceId 统计运行时设备
func (d *DeviceUtilizationDaoImpl) statDeviceStatByDeviceId(c *gin.Context, startTime string, endTime string,
	deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error) {
	return d.metricDao.StatDeviceStatusByDeviceId(c, startTime, endTime, deviceId, status)
}

// statDeviceProcessByDeviceId 统计设备运行过程
func (d *DeviceUtilizationDaoImpl) statDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string,
	deviceId int64, interval string,
	status int) (*devicemonitormodel.DeviceProcessList, error) {
	return d.metricDao.StatDeviceProcessByDeviceId(c, startTime, endTime, deviceId, interval, status)
}
