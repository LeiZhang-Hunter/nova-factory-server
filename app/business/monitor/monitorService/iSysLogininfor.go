package monitorService

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorModels"
)

type ILogininforService interface {
	SelectLogininforList(c *gin.Context, logininfor *monitorModels.LogininforDQL) (list []*monitorModels.Logininfor, total int64)
	ExportLogininfor(c *gin.Context, logininfor *monitorModels.LogininforDQL) (data []byte)
	InsertLogininfor(c context.Context, loginUser *monitorModels.Logininfor)
	DeleteLogininforByIds(c *gin.Context, infoIds []int64)
	CleanLogininfor(c *gin.Context)
}
