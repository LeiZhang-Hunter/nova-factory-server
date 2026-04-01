package monitordao

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
)

type IOperLog interface {
	InsertOperLog(ctx context.Context, operLog *monitormodels.SysOperLog)
	SelectOperLogList(ctx context.Context, openLog *monitormodels.SysOperLogDQL) (list []*monitormodels.SysOperLog, total int64)
	SelectOperLogListAll(ctx context.Context, openLog *monitormodels.SysOperLogDQL) (list []*monitormodels.SysOperLog)
	DeleteOperLogByIds(ctx context.Context, operIds []int64)
	SelectOperLogById(ctx context.Context, operId int64) (operLog *monitormodels.SysOperLog)
	CleanOperLog(ctx context.Context)
}
