package aiDataSetDaoImpl

import (
	"strings"
	"time"

	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAiConversationDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIAiConversationDaoImpl(db *gorm.DB) aidatasetdao.IAiConversationDao {
	return &IAiConversationDaoImpl{
		db:    db,
		table: "ai_conversations",
	}
}

func (i *IAiConversationDaoImpl) Create(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error) {
	data := &aidatasetmodels.AiConversation{
		Name:           req.Name,
		Message:        req.Message,
		LLMProviderID:  req.LLMProviderID,
		LLMModelID:     req.LLMModelID,
		EnableThinking: req.EnableThinking,
		ChatMode:       req.ChatMode,
		DeptID:         baizeContext.GetDeptId(c),
		State:          commonStatus.NORMAL,
	}
	data.SetCreateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (i *IAiConversationDaoImpl) Update(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error) {
	data := &aidatasetmodels.AiConversation{}
	if err := i.db.WithContext(c).Table(i.table).
		Where("id = ?", req.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(data).Error; err != nil {
		return nil, err
	}
	data.Name = req.Name
	data.Message = req.Message
	data.LLMProviderID = req.LLMProviderID
	data.LLMModelID = req.LLMModelID
	data.EnableThinking = req.EnableThinking
	data.ChatMode = req.ChatMode
	data.SetUpdateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).
		Where("id = ?", req.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"name":            data.Name,
			"message":         data.Message,
			"llm_provider_id": data.LLMProviderID,
			"llm_model_id":    data.LLMModelID,
			"enable_thinking": data.EnableThinking,
			"chat_mode":       data.ChatMode,
			"update_by":       data.UpdateBy,
			"update_time":     data.UpdateTime,
		}).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (i *IAiConversationDaoImpl) List(c *gin.Context, req *aidatasetmodels.AiConversationQuery) (*aidatasetmodels.AiConversationListData, error) {
	db := i.db.WithContext(c).Table(i.table).Where("state = ?", commonStatus.NORMAL).
		Where("dept_id = ?", baizeContext.GetDeptId(c))
	if strings.TrimSpace(req.Name) != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*aidatasetmodels.AiConversation, 0)
	if err := db.Order("id DESC").Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &aidatasetmodels.AiConversationListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (i *IAiConversationDaoImpl) Remove(c *gin.Context, ids []int64) error {
	return i.db.WithContext(c).Table(i.table).Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}
