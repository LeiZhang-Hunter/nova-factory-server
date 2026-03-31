package aiDataSetDaoImpl

import (
	"errors"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (a *AiLLMDaoImpl) ListByFactory(c *gin.Context, factory string) ([]*aiDataSetModels.AiLLMEntity, error) {
	rows := make([]*aiDataSetModels.AiLLMEntity, 0)
	if factory == "" {
		return rows, nil
	}
	if err := a.db.WithContext(c).Table(a.tableName).Where("state = ?", commonStatus.NORMAL).
		Where("fid = ?", factory).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (a *AiLLMDaoImpl) UpsertFactoryLLMs(tx *gorm.DB, providerName string, llms []*aiDataSetModels.FactoryLLMUpsert, status string, now time.Time) error {
	for _, llm := range llms {
		if llm == nil || llm.LLMName == "" {
			continue
		}
		toolFlag := 0
		if llm.IsTools {
			toolFlag = 1
		}
		var row aiDataSetModels.AiLLMEntity
		kt := time.Now()
		err := tx.Table(a.tableName).
			Where("fid = ?", providerName).Where("llm_name = ?", llm.LLMName).Take(&row).Error
		if err == nil {
			update := aiDataSetModels.AiLLMEntity{
				ModelType: llm.ModelType,
				MaxTokens: llm.MaxTokens,
				Tags:      llm.Tags,
				IsTools:   int32(toolFlag),
				Status:    status,
				State:     0,
				BaseEntity: baize.BaseEntity{
					CreateTime: &kt,
					UpdateTime: &kt,
				},
			}
			if err = tx.Table(a.tableName).Where("fid = ?", providerName).Where("llm_name = ?", llm.LLMName).
				Select("model_type", "max_tokens", "tags", "is_tools", "status", "state", "update_time_db").Updates(&update).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		create := aiDataSetModels.AiLLMEntity{
			FID:       providerName,
			LlmName:   llm.LLMName,
			ModelType: llm.ModelType,
			MaxTokens: llm.MaxTokens,
			Tags:      llm.Tags,
			IsTools:   int32(toolFlag),
			Status:    status,
			State:     0,
			BaseEntity: baize.BaseEntity{
				CreateTime: &kt,
				UpdateTime: &kt,
			},
		}
		if err = tx.Table(a.tableName).Create(&create).Error; err != nil {
			return err
		}
	}
	return nil
}

func (a *AiLLMDaoImpl) ListExistingLLMNames(c *gin.Context, llmNames []string) ([]string, error) {
	if len(llmNames) == 0 {
		return make([]string, 0), nil
	}
	type row struct {
		LlmName string `gorm:"column:llm_name"`
	}
	rows := make([]*row, 0)
	if err := a.db.WithContext(c).Table(a.tableName).Where("state = ?", commonStatus.NORMAL).
		Where("llm_name IN ?", llmNames).Find(&rows).Error; err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(rows))
	for _, r := range rows {
		ret = append(ret, r.LlmName)
	}
	return ret, nil
}
