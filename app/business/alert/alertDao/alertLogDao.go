package alertDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertLogDao interface {
	Export(c *gin.Context, data []*alertModels.SysAlertLog) error
	List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.SysAlertLogList, error)
	UpdateReason(c *gin.Context, id int64, reason string) error
}
