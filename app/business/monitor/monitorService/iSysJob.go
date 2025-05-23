package monitorService

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorModels"
)

type IJobService interface {
	SelectJobList(c *gin.Context, job *monitorModels.JobDQL) (list []*monitorModels.JobVo, total int64)
	SelectJobById(c *gin.Context, id int64) (job *monitorModels.JobVo)
	DeleteJobByIds(c *gin.Context, jobIds []int64)
	ChangeStatus(c *gin.Context, job *monitorModels.JobStatus)
	Run(c *gin.Context, job *monitorModels.JobStatus)
	InsertJob(c *gin.Context, job *monitorModels.JobDML)
	UpdateJob(c *gin.Context, job *monitorModels.JobDML)
	StartRunCron(c context.Context, jo *monitorModels.JobRedis)
	DeleteRunCron(c context.Context, jo *monitorModels.JobRedis)
	FunIsExist(invokeTarget string) bool
	GetFunList() []string
	SelectJobLogList(c *gin.Context, job *monitorModels.JobLogDql) (list []*monitorModels.JobLog, total int64)
	SelectJobLogById(c *gin.Context, id int64) (vo *monitorModels.JobLog)
	SelectJobIdAndNameAll(c *gin.Context) (list []*monitorModels.JobIdAndName)
}
