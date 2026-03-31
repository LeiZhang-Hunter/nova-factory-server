package alertservice

import (
	alertModels2 "nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertLogService interface {
	Export(c *gin.Context, data alertModels2.AlertLogData) error
	List(c *gin.Context, req *alertModels2.SysAlertLogListReq) (*alertModels2.NovaAlertLogList, error)
	Info(c *gin.Context, objectId uint64) (*alertModels2.NovaAlertLog, error)
	// Count 统计clickhouse数据库
	Count(c *gin.Context) (int64, error)
}
