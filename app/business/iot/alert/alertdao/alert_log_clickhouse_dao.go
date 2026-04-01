package alertdao

import (
	alertModels2 "nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertLogClickhouseDao interface {
	Export(c *gin.Context, data []*alertModels2.NovaAlertLog) error
	List(c *gin.Context, req *alertModels2.SysAlertLogListReq) (*alertModels2.NovaAlertLogList, error)
	GetByObjectId(c *gin.Context, objectId uint64) (*alertModels2.NovaAlertLog, error)
	// Count 统计clickhouse数据库
	Count(c *gin.Context) (int64, error)
}
