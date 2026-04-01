package monitorservice

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitormodels"

	"github.com/gin-gonic/gin"
)

type ILogininforService interface {
	SelectLogininforList(c *gin.Context, logininfor *monitormodels.LogininforDQL) (list []*monitormodels.Logininfor, total int64)
	ExportLogininfor(c *gin.Context, logininfor *monitormodels.LogininforDQL) (data []byte)
	InsertLogininfor(c context.Context, loginUser *monitormodels.Logininfor)
	DeleteLogininforByIds(c *gin.Context, infoIds []int64)
	CleanLogininfor(c *gin.Context)
}
