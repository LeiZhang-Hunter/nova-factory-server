package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"time"

	"gorm.io/gorm"
)

type IAiModelProviderDao interface {
	ListWithLLM(c *gin.Context, req *aiDataSetModels.SysAiModelProviderListReq) (*aiDataSetModels.SysAiModelProviderListData, error)
	UpsertFactoryProvider(tx *gorm.DB, item *aiDataSetModels.FactoryProviderUpsert, status int32, rank int32, now time.Time) (int64, error)
}
