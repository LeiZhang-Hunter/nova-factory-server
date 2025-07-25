package deviceMonitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type IDeviceDataReportService interface {
	DevList(c *gin.Context) ([]deviceMonitorModel.SysIotDbDevMapData, error)
}
