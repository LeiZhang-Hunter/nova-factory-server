package deviceMonitorServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
)

type IDeviceDataReportServiceImpl struct {
	dao deviceMonitorDao.IDeviceDataReportDao
}

func NewIDeviceDataReportServiceImpl(dao deviceMonitorDao.IDeviceDataReportDao) deviceMonitorService.IDeviceDataReportService {
	return &IDeviceDataReportServiceImpl{
		dao: dao,
	}
}

func (i *IDeviceDataReportServiceImpl) DevList(c *gin.Context) ([]deviceMonitorModel.SysIotDbDevMapData, error) {
	list, err := i.dao.DevList(c)
	if err != nil {
		return nil, err
	}
	return list, nil
}
