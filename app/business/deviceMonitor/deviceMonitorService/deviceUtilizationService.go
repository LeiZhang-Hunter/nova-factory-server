package deviceMonitorService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
)

type DeviceUtilizationService interface {
	Stat(c *gin.Context, req *deviceMonitorModel.DeviceUtilizationReq)
}
