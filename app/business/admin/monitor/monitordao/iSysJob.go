package monitordao

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
)

type IJobDao interface {
	SelectJobList(ctx context.Context, job *monitormodels.JobDQL) (list []*monitormodels.JobVo, total int64)
	SelectJobAll(ctx context.Context) (list []*monitormodels.JobVo)
	SelectJobById(ctx context.Context, id int64) (job *monitormodels.JobVo)
	SelectJobByInvokeTarget(ctx context.Context, invokeTarget string) (job *monitormodels.JobVo)
	DeleteJobById(ctx context.Context, id int64)
	UpdateJob(ctx context.Context, job *monitormodels.JobDML)
	InsertJob(ctx context.Context, job *monitormodels.JobDML)
	DeleteJobByIds(ctx context.Context, id []int64)
	InsertJobLog(ctx context.Context, job *monitormodels.JobLog)
	SelectJobLogList(ctx context.Context, job *monitormodels.JobLogDql) (list []*monitormodels.JobLog, total int64)
	SelectJobLogById(ctx context.Context, id int64) (vo *monitormodels.JobLog)
	SelectJobIdAndNameAll(ctx context.Context) (list []*monitormodels.JobIdAndName)
}
