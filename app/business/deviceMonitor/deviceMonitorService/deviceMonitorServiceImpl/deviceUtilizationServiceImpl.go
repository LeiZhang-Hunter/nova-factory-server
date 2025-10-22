package deviceMonitorServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
)

type DeviceUtilizationServiceImpl struct {
	deviceUtilizationDao deviceMonitorDao.DeviceUtilizationDao
}

func NewDeviceUtilizationServiceImpl(deviceUtilizationDao deviceMonitorDao.DeviceUtilizationDao) deviceMonitorService.DeviceUtilizationService {
	return &DeviceUtilizationServiceImpl{
		deviceUtilizationDao: deviceUtilizationDao,
	}
}

func (d *DeviceUtilizationServiceImpl) Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) {
	//runMetrics, err := d.metricCDao.Query(c, &metricModels.MetricDataQueryReq{
	//ByDevice: true,
	//})
	//if err != nil {
	//	return
	//}
	d.deviceUtilizationDao.Stat(c, req)
}
