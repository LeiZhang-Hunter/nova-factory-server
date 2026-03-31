package monitorServiceImpl

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitordao"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/business/admin/monitor/monitorservice"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/snowflake"
	"time"

	"github.com/gin-gonic/gin"
)

type OperLogService struct {
	old monitordao.IOperLog
}

func NewOperLog(old monitordao.IOperLog) monitorservice.ISysOperLogService {
	return &OperLogService{old: old}
}

func (ols *OperLogService) InsertOperLog(c context.Context, operLog *monitormodels.SysOperLog) {
	operLog.OperId = snowflake.GenID()
	operLog.OperTime = time.Now()
	ols.old.InsertOperLog(c, operLog)
}
func (ols *OperLogService) SelectOperLogList(c *gin.Context, openLog *monitormodels.SysOperLogDQL) (list []*monitormodels.SysOperLog, total int64) {
	list, total = ols.old.SelectOperLogList(c, openLog)
	return

}
func (ols *OperLogService) ExportOperLog(c *gin.Context, openLog *monitormodels.SysOperLogDQL) (data []byte) {
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
func (ols *OperLogService) SelectOperLogById(c *gin.Context, operId int64) (operLog *monitormodels.SysOperLog) {
	return ols.old.SelectOperLogById(c, operId)
}
func (ols *OperLogService) CleanOperLog(c *gin.Context) {
	ols.old.CleanOperLog(c)
}
