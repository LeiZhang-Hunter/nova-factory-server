package aiDataSetDao

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type IAiLLMDao interface {
	ListByFIDs(c *gin.Context, fids []string) ([]*aiDataSetModels.SysAiLLM, error)
	ListByFactory(c *gin.Context, factory string) ([]*aiDataSetModels.AiLLMEntity, error)
	UpsertFactoryLLMs(tx *gorm.DB, providerName string, llms []*aiDataSetModels.FactoryLLMUpsert, status string, now time.Time) error
	ListExistingLLMNames(c *gin.Context, llmNames []string) ([]string, error)
}
