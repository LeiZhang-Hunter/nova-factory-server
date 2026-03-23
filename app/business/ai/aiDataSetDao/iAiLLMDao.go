package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IAiLLMDao interface {
	ListByFIDs(c *gin.Context, fids []string) ([]*aiDataSetModels.SysAiLLM, error)
}
