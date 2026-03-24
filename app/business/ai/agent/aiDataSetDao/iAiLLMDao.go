package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"time"

	"gorm.io/gorm"
)

type IAiLLMDao interface {
	ListByFIDs(c *gin.Context, fids []string) ([]*aiDataSetModels.SysAiLLM, error)
	UpsertFactoryLLMs(tx *gorm.DB, providerName string, llms []*aiDataSetModels.FactoryLLMUpsert, status string, now time.Time) error
}
