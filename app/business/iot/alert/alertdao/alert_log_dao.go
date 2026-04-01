package alertdao

import (
	"nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertLogDao interface {
	Export(c *gin.Context, data []*alertmodels.SysAlertLog) error
	List(c *gin.Context, req *alertmodels.SysAlertLogListReq) (*alertmodels.SysAlertLogList, error)
	UpdateReason(c *gin.Context, id int64, reason string) error
}
