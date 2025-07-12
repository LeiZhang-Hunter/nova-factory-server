package deviceMonitorDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type IDeviceDataReportDao interface {
	DevList(c *gin.Context) ([]deviceMonitorModel.SysIotDbDevMapData, error)
	GetDevList(c *gin.Context, dev []string) ([]deviceMonitorModel.SysIotDbDevMapData, error)
}
