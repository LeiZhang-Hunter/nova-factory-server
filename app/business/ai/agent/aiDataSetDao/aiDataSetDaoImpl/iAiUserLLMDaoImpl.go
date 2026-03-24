package aiDataSetDaoImpl

import (
	"errors"
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAiUserLLMDaoImpl struct {
	db     *gorm.DB
	table  string
	llmDao aiDataSetDao.IAiLLMDao
}

func NewIAiUserLLMDaoImpl(db *gorm.DB, llmDao aiDataSetDao.IAiLLMDao) aiDataSetDao.IAiUserLLMDao {
	return &IAiUserLLMDaoImpl{
		db:     db,
		table:  "ai_user_llm",
		llmDao: llmDao,
	}
}

func (i *IAiUserLLMDaoImpl) Set(c *gin.Context, req *aiDataSetModels.SetSysUserLLM) (*aiDataSetModels.SysUserLLM, error) {
	if req.LLMFactory == "" {
		return nil, errors.New("llm_factory不能为空")
	}
	userID := 0
	llms, err := i.llmDao.ListByFactory(c, req.LLMFactory)
	if err != nil {
		return nil, err
	}
	if len(llms) == 0 {
		return nil, errors.New("该工厂下未找到模型")
	}
	var result *aiDataSetModels.SysUserLLM
	err = i.db.Transaction(func(tx *gorm.DB) error {
		for _, llm := range llms {
			data := &aiDataSetModels.SysUserLLM{
				UserID:     int64(userID),
				LLMFactory: req.LLMFactory,
				ModelType:  llm.ModelType,
				LLMName:    llm.LlmName,
				APIKey:     req.APIKey,
				APIBase:    req.APIBase,
				MaxTokens:  llm.MaxTokens,
				UsedTokens: 0,
				Status:     llm.Status,
			}
			var exists aiDataSetModels.SysUserLLM
			err := tx.Table(i.table).Where("user_id = ?", userID).Where("llm_factory = ?", req.LLMFactory).
				Where("llm_name = ?", llm.LlmName).First(&exists).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = tx.Table(i.table).Create(data).Error; err != nil {
					return err
				}
				result = data
				continue
			}
			if err != nil {
				return err
			}
			if err = tx.Table(i.table).Where("user_id = ?", userID).Where("llm_factory = ?", req.LLMFactory).
				Where("llm_name = ?", llm.LlmName).
				Select("model_type", "api_key", "api_base", "max_tokens", "status").Updates(data).Error; err != nil {
				return err
			}
			if err = tx.Table(i.table).Where("user_id = ?", userID).Where("llm_factory = ?", req.LLMFactory).
				Where("llm_name = ?", llm.LlmName).First(&exists).Error; err != nil {
				return err
			}
			result = &exists
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		if err = i.db.Table(i.table).Where("user_id = ?", userID).Where("llm_factory = ?", req.LLMFactory).First(&aiDataSetModels.SysUserLLM{}).Error; err != nil {
			return nil, err
		}
		return nil, errors.New("写入失败")
	}
	return result, nil
}

func (i *IAiUserLLMDaoImpl) Get(c *gin.Context, req *aiDataSetModels.GetSysUserLLMReq) ([]*aiDataSetModels.SysUserLLM, error) {
	userID := int64(0)
	rows := make([]*aiDataSetModels.SysUserLLM, 0)
	db := i.db.Table(i.table).Where("user_id = ?", userID)
	if req != nil && req.LLMFactory != "" {
		db = db.Where("llm_factory = ?", req.LLMFactory)
	}
	if err := db.Order("llm_factory asc,llm_name asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
