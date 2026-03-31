package metricservice

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"

	"github.com/gin-gonic/gin"
)

type IDevMapService interface {
	GetDevList(c *gin.Context, dev []string) ([]devicemonitormodel.SysIotDbDevMapData, error)
}
