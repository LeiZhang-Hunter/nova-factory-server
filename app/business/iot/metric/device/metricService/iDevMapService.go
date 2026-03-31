package metricService

import (
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorModel"

	"github.com/gin-gonic/gin"
)

type IDevMapService interface {
	GetDevList(c *gin.Context, dev []string) ([]deviceMonitorModel.SysIotDbDevMapData, error)
}
