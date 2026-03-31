package monitorservice

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitormodels"

	"github.com/gin-gonic/gin"
)

type IJobService interface {
	SelectJobList(c *gin.Context, job *monitormodels.JobDQL) (list []*monitormodels.JobVo, total int64)
	SelectJobById(c *gin.Context, id int64) (job *monitormodels.JobVo)
	DeleteJobByIds(c *gin.Context, jobIds []int64)
	ChangeStatus(c *gin.Context, job *monitormodels.JobStatus)
	Run(c *gin.Context, job *monitormodels.JobStatus)
	InsertJob(c *gin.Context, job *monitormodels.JobDML)
	UpdateJob(c *gin.Context, job *monitormodels.JobDML)
	StartRunCron(c context.Context, jo *monitormodels.JobRedis)
	DeleteRunCron(c context.Context, jo *monitormodels.JobRedis)
	FunIsExist(invokeTarget string) bool
	GetFunList() []string
	SelectJobLogList(c *gin.Context, job *monitormodels.JobLogDql) (list []*monitormodels.JobLog, total int64)
	SelectJobLogById(c *gin.Context, id int64) (vo *monitormodels.JobLog)
	SelectJobIdAndNameAll(c *gin.Context) (list []*monitormodels.JobIdAndName)
}
