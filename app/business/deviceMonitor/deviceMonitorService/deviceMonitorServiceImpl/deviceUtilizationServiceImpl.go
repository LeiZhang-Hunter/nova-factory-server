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

func (d *DeviceUtilizationServiceImpl) Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) ([]*deviceMonitorModel.DeviceUtilizationData, error) {
	//runMetrics, err := d.metricCDao.Query(c, &metricModels.MetricDataQueryReq{
	//ByDevice: true,
	//})
	//if err != nil {
	//	return
	//}
	return d.deviceUtilizationDao.Stat(c, req)
}

func (d *DeviceUtilizationServiceImpl) Search(c *gin.Context,
	req *deviceMonitorModel.DeviceUtilizationReq) (*deviceMonitorModel.DeviceUtilizationPublicDataList, error) {
	return d.deviceUtilizationDao.Search(c, req)
}
