package monitorservice

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitormodels"

	"github.com/gin-gonic/gin"
)

type ISysOperLogService interface {
	InsertOperLog(c context.Context, operLog *monitormodels.SysOperLog)
	SelectOperLogList(c *gin.Context, openLog *monitormodels.SysOperLogDQL) (list []*monitormodels.SysOperLog, total int64)
	ExportOperLog(c *gin.Context, logininfor *monitormodels.SysOperLogDQL) (data []byte)
	DeleteOperLogByIds(c *gin.Context, operIds []int64)
	SelectOperLogById(c *gin.Context, operId int64) (operLogList *monitormodels.SysOperLog)
	CleanOperLog(c *gin.Context)
}
