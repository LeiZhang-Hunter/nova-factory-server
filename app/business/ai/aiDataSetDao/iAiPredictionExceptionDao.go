package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IAiPredictionExceptionDao interface {
	Set(c *gin.Context, data *aiDataSetModels.SetSysAiPredictionException) (*aiDataSetModels.SysAiPredictionException, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aiDataSetModels.SysAiPredictionExceptionListReq) (*aiDataSetModels.SysAiPredictionExceptionList, error)
	All(c *gin.Context) ([]*aiDataSetModels.SysAiPredictionException, error)
}
