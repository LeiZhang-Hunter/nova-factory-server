package deviceMonitorServiceImpl

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitordao"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitorservice"

	"github.com/gin-gonic/gin"
)

type DeviceUtilizationServiceImpl struct {
	deviceUtilizationDao devicemonitordao.DeviceUtilizationDao
}

func NewDeviceUtilizationServiceImpl(deviceUtilizationDao devicemonitordao.DeviceUtilizationDao) devicemonitorservice.DeviceUtilizationService {
	return &DeviceUtilizationServiceImpl{
		deviceUtilizationDao: deviceUtilizationDao,
	}
}

func (d *DeviceUtilizationServiceImpl) Stat(c *gin.Context, req *devicemonitormodel.DeviceUtilizationReq) ([]*devicemonitormodel.DeviceUtilizationData, error) {
	//runMetrics, err := d.metricCDao.Query(c, &metricmodels.MetricDataQueryReq{
	//ByDevice: true,
	//})
	//if err != nil {
	//	return
	//}
	return d.deviceUtilizationDao.Stat(c, req)
}

func (d *DeviceUtilizationServiceImpl) Search(c *gin.Context,
	req *devicemonitormodel.DeviceUtilizationReq) (*devicemonitormodel.DeviceUtilizationPublicDataList, error) {
	return d.deviceUtilizationDao.Search(c, req)
}

func (d *DeviceUtilizationServiceImpl) SearchV2(c *gin.Context, req *devicemonitormodel.DeviceUtilizationReq) (*devicemonitormodel.DeviceUtilizationPublicDataListV2, error) {
	return d.deviceUtilizationDao.SearchV2(c, req)
}

func (d *DeviceUtilizationServiceImpl) GetDeviceUtilization(c *gin.Context, req *devicemonitormodel.DeviceUtilizationReq) (*devicemonitormodel.DeviceRunProcess, error) {
	return d.deviceUtilizationDao.GetDeviceUtilization(c, req)
}
