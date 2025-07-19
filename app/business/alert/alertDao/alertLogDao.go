package alertDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertLogDao interface {
	Export(c *gin.Context, data []*alertModels.SysAlertLog) error
}
