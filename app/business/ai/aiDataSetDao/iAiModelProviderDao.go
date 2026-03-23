package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IAiModelProviderDao interface {
	ListWithLLM(c *gin.Context, req *aiDataSetModels.SysAiModelProviderListReq) (*aiDataSetModels.SysAiModelProviderListData, error)
}
