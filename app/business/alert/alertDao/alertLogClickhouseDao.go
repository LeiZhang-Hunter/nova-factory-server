package alertDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertLogClickhouseDao interface {
	Export(c *gin.Context, data []*alertModels.NovaAlertLog) error
	List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.NovaAlertLogList, error)
}
