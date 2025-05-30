package monitorServiceImpl

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorDao"
	"nova-factory-server/app/business/monitor/monitorModels"
	"nova-factory-server/app/business/monitor/monitorService"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/snowflake"
	"time"
)

type OperLogService struct {
	old monitorDao.IOperLog
}

func NewOperLog(old monitorDao.IOperLog) monitorService.ISysOperLogService {
	return &OperLogService{old: old}
}

func (ols *OperLogService) InsertOperLog(c context.Context, operLog *monitorModels.SysOperLog) {
	operLog.OperId = snowflake.GenID()
	operLog.OperTime = time.Now()
	ols.old.InsertOperLog(c, operLog)
}
func (ols *OperLogService) SelectOperLogList(c *gin.Context, openLog *monitorModels.SysOperLogDQL) (list []*monitorModels.SysOperLog, total int64) {
	list, total = ols.old.SelectOperLogList(c, openLog)
	return

}
func (ols *OperLogService) ExportOperLog(c *gin.Context, openLog *monitorModels.SysOperLogDQL) (data []byte) {
	list := ols.old.SelectOperLogListAll(c, openLog)
	toExcel, err := excel.SliceToExcel(list)
	if err != nil {
		panic(err)
	}
	buffer, err := toExcel.WriteToBuffer()
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func (ols *OperLogService) DeleteOperLogByIds(c *gin.Context, operIds []int64) {
	ols.old.DeleteOperLogByIds(c, operIds)
}
func (ols *OperLogService) SelectOperLogById(c *gin.Context, operId int64) (operLog *monitorModels.SysOperLog) {
	return ols.old.SelectOperLogById(c, operId)
}
func (ols *OperLogService) CleanOperLog(c *gin.Context) {
	ols.old.CleanOperLog(c)
}
