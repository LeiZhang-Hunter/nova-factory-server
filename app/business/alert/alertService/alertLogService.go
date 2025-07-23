package alertService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertLogService interface {
	Export(c *gin.Context, data alertModels.AlertLogData) error
	List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.SysAlertLogList, error)
}
