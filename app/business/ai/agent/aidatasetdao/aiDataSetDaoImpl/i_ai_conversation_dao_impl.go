package aiDataSetDaoImpl

import (
	"errors"
	"nova-factory-server/app/utils/snowflake"
	"strings"

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
	if req.ChatMode == "" {
		req.ChatMode = "task"
	}
	if req.EnableThinking == nil {
		ret := true
		req.EnableThinking = &ret
	}
	data := &aidatasetmodels.AiConversation{
		Name:           req.Name,
		AgentID:        req.AgentID,
		AgentType:      req.AgentType,
		Message:        req.Message,
		LLMProviderID:  req.LLMProviderID,
		LLMModelID:     req.LLMModelID,
		EnableThinking: req.EnableThinking,
		ChatMode:       req.ChatMode,
		//DeptID:         baizeContext.GetDeptId(c),
		State: commonStatus.NORMAL,
	}
	data.ID = snowflake.GenID()
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
	data.AgentID = req.AgentID
	data.AgentType = req.AgentType
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
			"agent_id":        data.AgentID,
			"agent_type":      data.AgentType,
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

func (i *IAiConversationDaoImpl) GetByID(c *gin.Context, id int64) (*aidatasetmodels.AiConversation, error) {
	data := &aidatasetmodels.AiConversation{}
	if err := i.db.WithContext(c).Table(i.table).
		Where("id = ?", id).
		Where("create_by = ?", baizeContext.GetUserId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (i *IAiConversationDaoImpl) List(c *gin.Context, req *aidatasetmodels.AiConversationQuery) (*aidatasetmodels.AiConversationListData, error) {
	db := i.db.WithContext(c).Table(i.table)
	if req.ID != 0 {
		db = db.Where("id = ?", req.ID)
	}
	if strings.TrimSpace(req.Name) != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	db.Where("create_by = ?", baizeContext.GetUserId(c))
	var total int64
	db = db.Where("state = ?", commonStatus.NORMAL)
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
	if len(ids) == 0 {
		return nil
	}
	return i.db.WithContext(c).Table(i.table).
		Where("id IN ?", ids).
		Where("create_by = ?", baizeContext.GetUserId(c)).
		Delete(&aidatasetmodels.AiConversation{}).Error
}
