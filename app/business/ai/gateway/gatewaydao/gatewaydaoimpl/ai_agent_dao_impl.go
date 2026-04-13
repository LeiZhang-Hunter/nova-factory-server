package gatewaydaoimpl

import (
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIAgentDaoImpl 提供智能体配置的数据访问能力。
type AIAgentDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewAIAgentDao 创建智能体配置DAO。
func NewAIAgentDao(db *gorm.DB) gatewaydao.IAIAgentDao {
	return &AIAgentDaoImpl{
		db:    db,
		table: "ai_agents",
	}
}

// Create 新增智能体配置。
func (a *AIAgentDaoImpl) Create(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error) {
	item := buildAIAgentModel(c, req)
	item.ID = snowflake.GenID()
	item.State = commonStatus.NORMAL
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Update 修改智能体配置。
func (a *AIAgentDaoImpl) Update(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error) {
	item := buildAIAgentModel(c, req)
	item.ID = req.ID
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).
		Where("id = ?", item.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Select("name", "prompt", "default_llm_provider_id", "default_llm_model_id", "llm_temperature", "llm_top_p",
			"llm_max_tokens", "enable_llm_temperature", "enable_llm_top_p", "enable_llm_max_tokens",
			"llm_max_context_count", "retrieval_top_k", "retrieval_match_threshold", "sandbox_mode", "sandbox_network",
			"work_dir", "mcp_enabled", "mcp_server_ids", "mcp_server_enabled_ids", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return a.GetByID(c, item.ID)
}

// DeleteByIDs 删除智能体配置。
func (a *AIAgentDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return a.db.WithContext(c).Table(a.table).
		Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

// GetByID 查询智能体配置详情。
func (a *AIAgentDaoImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgent, error) {
	var item gatewaymodels.AIAgent
	if err := a.db.WithContext(c).Table(a.table).
		Where("id = ?", id).
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

// List 查询智能体配置列表。
func (a *AIAgentDaoImpl) List(c *gin.Context, req *gatewaymodels.AIAgentQuery) (*gatewaymodels.AIAgentListData, error) {
	db := a.db.WithContext(c).Table(a.table).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if name := strings.TrimSpace(req.Name); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if sandboxMode := strings.TrimSpace(req.SandboxMode); sandboxMode != "" {
		db = db.Where("sandbox_mode = ?", sandboxMode)
	}
	if req.MCPEnabled != nil {
		db = db.Where("mcp_enabled = ?", req.MCPEnabled)
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
	rows := make([]*gatewaymodels.AIAgent, 0)
	if err := db.Order("create_time DESC, id DESC").
		Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &gatewaymodels.AIAgentListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func buildAIAgentModel(c *gin.Context, req *gatewaymodels.AIAgentUpsert) *gatewaymodels.AIAgent {
	return &gatewaymodels.AIAgent{
		Name:                    req.Name,
		Prompt:                  req.Prompt,
		DefaultLLMProviderID:    req.DefaultLLMProviderID,
		DefaultLLMModelID:       req.DefaultLLMModelID,
		LLMTemperature:          req.LLMTemperature,
		LLMTopP:                 req.LLMTopP,
		LLMMaxTokens:            req.LLMMaxTokens,
		EnableLLMTemperature:    req.EnableLLMTemperature,
		EnableLLMTopP:           req.EnableLLMTopP,
		EnableLLMMaxTokens:      req.EnableLLMMaxTokens,
		LLMMaxContextCount:      req.LLMMaxContextCount,
		RetrievalTopK:           req.RetrievalTopK,
		RetrievalMatchThreshold: req.RetrievalMatchThreshold,
		SandboxMode:             req.SandboxMode,
		SandboxNetwork:          req.SandboxNetwork,
		WorkDir:                 req.WorkDir,
		MCPEnabled:              req.MCPEnabled,
		MCPServerIDs:            req.MCPServerIDs,
		MCPServerEnabledIDs:     req.MCPServerEnabledIDs,
		DeptID:                  baizeContext.GetDeptId(c),
	}
}
