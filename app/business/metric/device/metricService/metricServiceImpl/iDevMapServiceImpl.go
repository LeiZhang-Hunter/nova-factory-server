package metricServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricService"
)

type IDevMapServiceImpl struct {
	mapDao deviceMonitorDao.IDeviceDataReportDao
}

func NewIDevMapServiceImpl(mapDao deviceMonitorDao.IDeviceDataReportDao) metricService.IDevMapService {
	return &IDevMapServiceImpl{
		mapDao: mapDao,
	}
}

func (i *IDevMapServiceImpl) GetDevList(c *gin.Context, dev []string) ([]deviceMonitorModel.SysIotDbDevMapData, error) {
	return i.mapDao.GetDevList(c, dev)
}
