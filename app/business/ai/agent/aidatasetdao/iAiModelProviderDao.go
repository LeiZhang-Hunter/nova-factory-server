package aidatasetdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"time"

	"gorm.io/gorm"
)

type IAiModelProviderDao interface {
	ListWithLLM(c *gin.Context, req *aidatasetmodels.SysAiModelProviderListReq) (*aidatasetmodels.SysAiModelProviderListData, error)
	UpsertFactoryProvider(tx *gorm.DB, item *aidatasetmodels.FactoryProviderUpsert, status int32, rank int32, now time.Time) (int64, error)
}
