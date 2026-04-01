package aidatasetdao

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type IAiLLMDao interface {
	ListByFIDs(c *gin.Context, fids []string) ([]*aidatasetmodels.SysAiLLM, error)
	ListByFactory(c *gin.Context, factory string) ([]*aidatasetmodels.AiLLMEntity, error)
	UpsertFactoryLLMs(tx *gorm.DB, providerName string, llms []*aidatasetmodels.FactoryLLMUpsert, status string, now time.Time) error
	ListExistingLLMNames(c *gin.Context, llmNames []string) ([]string, error)
}
