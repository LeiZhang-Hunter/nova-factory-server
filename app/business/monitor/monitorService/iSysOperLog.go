package monitorService

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorModels"
)

type ISysOperLogService interface {
	InsertOperLog(c context.Context, operLog *monitorModels.SysOperLog)
	SelectOperLogList(c *gin.Context, openLog *monitorModels.SysOperLogDQL) (list []*monitorModels.SysOperLog, total int64)
	ExportOperLog(c *gin.Context, logininfor *monitorModels.SysOperLogDQL) (data []byte)
	DeleteOperLogByIds(c *gin.Context, operIds []int64)
	SelectOperLogById(c *gin.Context, operId int64) (operLogList *monitorModels.SysOperLog)
	CleanOperLog(c *gin.Context)
}
