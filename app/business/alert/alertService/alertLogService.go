package alertService

import (
	"nova-factory-server/app/business/alert/alertModels"

	"github.com/gin-gonic/gin"
)

type AlertLogService interface {
	Export(c *gin.Context, data alertModels.AlertLogData) error
	List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.NovaAlertLogList, error)
	Info(c *gin.Context, objectId uint64) (*alertModels.NovaAlertLog, error)
	// Count 统计clickhouse数据库
	Count(c *gin.Context) (int64, error)
}
