package aiDataSetDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
)

type AiLLMDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewAiLLMDaoImpl(db *gorm.DB) aiDataSetDao.IAiLLMDao {
	return &AiLLMDaoImpl{
		db:        db,
		tableName: "ai_llm",
	}
}

func (a *AiLLMDaoImpl) ListByFIDs(c *gin.Context, fids []string) ([]*aiDataSetModels.SysAiLLM, error) {
	llms := make([]*aiDataSetModels.SysAiLLM, 0)
	if len(fids) == 0 {
		return llms, nil
	}
	if err := a.db.WithContext(c).Table(a.tableName).Where("state = ?", commonStatus.NORMAL).
		Where("fid IN ?", fids).Order("llm_name ASC").Find(&llms).Error; err != nil {
		return nil, err
	}
	return llms, nil
}
