package gatewaydaoimpl

import (
	"errors"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AIAgentMessageDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAIAgentMessageDao(db *gorm.DB) gatewaydao.IAIAgentMessageDao {
	return &AIAgentMessageDaoImpl{
		db:    db,
		table: "ai_agent_messages",
	}
}

func (a *AIAgentMessageDaoImpl) Create(c *gin.Context, req *gatewaymodels.AIAgentMessageUpsert) (*gatewaymodels.AIAgentMessage, error) {
	item := &gatewaymodels.AIAgentMessage{
		ID:              snowflake.GenID(),
		ConversationID:  req.ConversationID,
		Role:            req.Role,
		Content:         req.Content,
		ProviderID:      req.ProviderID,
		ModelID:         req.ModelID,
		Status:          req.Status,
		Error:           req.Error,
		InputTokens:     req.InputTokens,
		OutputTokens:    req.OutputTokens,
		FinishReason:    req.FinishReason,
		ToolCalls:       req.ToolCalls,
		ToolCallID:      req.ToolCallID,
		ToolCallName:    req.ToolCallName,
		ThinkingContent: req.ThinkingContent,
		Segments:        req.Segments,
		ImagesJSON:      req.ImagesJSON,
		DeptID:          baizeContext.GetDeptId(c),
		State:           commonStatus.NORMAL,
	}
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (a *AIAgentMessageDaoImpl) Update(c *gin.Context, req *gatewaymodels.AIAgentMessageUpsert) (*gatewaymodels.AIAgentMessage, error) {
	item := &gatewaymodels.AIAgentMessage{
		ID:              req.ID,
		ConversationID:  req.ConversationID,
		Role:            req.Role,
		Content:         req.Content,
		ProviderID:      req.ProviderID,
		ModelID:         req.ModelID,
		Status:          req.Status,
		Error:           req.Error,
		InputTokens:     req.InputTokens,
		OutputTokens:    req.OutputTokens,
		FinishReason:    req.FinishReason,
		ToolCalls:       req.ToolCalls,
		ToolCallID:      req.ToolCallID,
		ToolCallName:    req.ToolCallName,
		ThinkingContent: req.ThinkingContent,
		Segments:        req.Segments,
		ImagesJSON:      req.ImagesJSON,
	}
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", item.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Select(
			"conversation_id",
			"role",
			"content",
			"provider_id",
			"model_id",
			"status",
			"error",
			"input_tokens",
			"output_tokens",
			"finish_reason",
			"tool_calls",
			"tool_call_id",
			"tool_call_name",
			"thinking_content",
			"segments",
			"images_json",
			"update_by",
			"update_time",
		).
		Updates(item).Error; err != nil {
		return nil, err
	}
	return a.GetByID(c, item.ID)
}

func (a *AIAgentMessageDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return a.db.WithContext(c).Table(a.table).Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

func (a *AIAgentMessageDaoImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgentMessage, error) {
	var item gatewaymodels.AIAgentMessage
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (a *AIAgentMessageDaoImpl) List(c *gin.Context, req *gatewaymodels.AIAgentMessageQuery) (*gatewaymodels.AIAgentMessageListData, error) {
	db := a.db.WithContext(c).Table(a.table)
	if req.ConversationID != 0 {
		db = db.Where("conversation_id = ?", req.ConversationID)
	}
	if len(req.ConversationIDs) > 0 {
		db = db.Where("conversation_id IN ?", req.ConversationIDs)
	}
	if req.Role != "" {
		db = db.Where("role = ?", req.Role)
	}
	if req.ProviderID != "" {
		db = db.Where("provider_id = ?", req.ProviderID)
	}
	if req.ModelID != "" {
		db = db.Where("model_id = ?", req.ModelID)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	db = db.Where("state = ?", commonStatus.NORMAL)
	rows := make([]*gatewaymodels.AIAgentMessage, 0)
	query := db.Order("conversation_id ASC").Order("id DESC")
	if req.Size > 0 {
		query = query.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size))
	}
	if err := query.Debug().Find(&rows).Error; err != nil {
		return nil, err
	}
	return &gatewaymodels.AIAgentMessageListData{
		Rows:  rows,
		Total: total,
	}, nil
}
