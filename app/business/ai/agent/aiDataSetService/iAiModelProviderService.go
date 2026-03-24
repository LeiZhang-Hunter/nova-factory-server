package aiDataSetService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
)

type IAiModelProviderService interface {
	ListWithLLM(c *gin.Context, req *aiDataSetModels.SysAiModelProviderListReq) (*aiDataSetModels.SysAiModelProviderListData, error)
}
