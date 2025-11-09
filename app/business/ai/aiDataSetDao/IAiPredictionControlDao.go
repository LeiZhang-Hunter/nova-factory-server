package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IAiPredictionControlDao interface {
	Set(c *gin.Context, data *aiDataSetModels.SetSysAiPredictionControl) (*aiDataSetModels.SysAiPredictionControl, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *aiDataSetModels.SysAiPredictionControlListReq) (*aiDataSetModels.SysAiPredictionControlList, error)
	Find(c *gin.Context) (*aiDataSetModels.SysAiPredictionControl, error)
}
