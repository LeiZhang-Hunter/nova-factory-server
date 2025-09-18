package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IAiPredictionListDao interface {
	Set(c *gin.Context, data *aiDataSetModels.SetSysAiPrediction) (*aiDataSetModels.SysAiPrediction, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aiDataSetModels.SysAiPredictionListReq) (*aiDataSetModels.SysAiPredictionList, error)
}
