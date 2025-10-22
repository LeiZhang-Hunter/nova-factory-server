package deviceMonitorDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type DeviceUtilizationDao interface {
	Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq) ([]*deviceMonitorModel.DeviceUtilizationData, error)
}
