package gatewaydaoimpl

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
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

// AISubAgentDaoImpl 提供子智能体配置的数据访问能力。
type AISubAgentDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewAISubAgentDao 创建子智能体配置DAO。
func NewAISubAgentDao(db *gorm.DB) gatewaydao.IAISubAgentDao {
	return &AISubAgentDaoImpl{
		db:    db,
		table: "ai_sub_agents",
	}
}

// Create 新增子智能体配置。
func (a *AISubAgentDaoImpl) Create(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error) {
	item := buildAISubAgentModel(c, req)
	item.ID = snowflake.GenID()
	item.State = commonStatus.NORMAL
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Update 修改子智能体配置。
func (a *AISubAgentDaoImpl) Update(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error) {
	item := buildAISubAgentModel(c, req)
	item.ID = req.ID
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).
		Where("id = ?", item.ID).
		Where("state = ?", commonStatus.NORMAL).
		Select("name", "type", "description", "instruction", "mcp_enabled", "mcp_server_ids",
			"mcp_server_enabled_ids", "local_tool_enabled", "local_tools", "allow_mcp_server_ids_tools",
			"enable", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return a.GetByID(c, item.ID)
}

// DeleteByIDs 删除子智能体配置。
func (a *AISubAgentDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return a.db.WithContext(c).Table(a.table).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

// GetByID 查询子智能体配置详情。
func (a *AISubAgentDaoImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AISubAgent, error) {
	var item gatewaymodels.AISubAgent
	if err := a.db.WithContext(c).Table(a.table).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	decodeSubAgentAllowMcpServerIdsTools(&item)
	return &item, nil
}

// List 查询子智能体配置列表。
func (a *AISubAgentDaoImpl) List(c *gin.Context, req *gatewaymodels.AISubAgentQuery) (*gatewaymodels.AISubAgentListData, error) {
	db := a.db.WithContext(c).Table(a.table).
		Where("state = ?", commonStatus.NORMAL)
	if name := strings.TrimSpace(req.Name); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if subAgentType := strings.TrimSpace(req.Type); subAgentType != "" {
		db = db.Where("type = ?", subAgentType)
	}
	if req.MCPEnabled != nil {
		db = db.Where("mcp_enabled = ?", req.MCPEnabled)
	}
	if req.LocalToolEnabled != nil {
		db = db.Where("local_tool_enabled = ?", req.LocalToolEnabled)
	}
	if req.Enable != nil {
		db = db.Where("enable = ?", req.Enable)
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
	rows := make([]*gatewaymodels.AISubAgent, 0)
	if err := db.Order("create_time DESC, id DESC").
		Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		decodeSubAgentAllowMcpServerIdsTools(row)
	}
	return &gatewaymodels.AISubAgentListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func buildAISubAgentModel(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) *gatewaymodels.AISubAgent {
	mCPServerIDs, err := json.Marshal(req.MCPServerIDs)
	if err != nil {
		zap.L().Error("buildAISubAgentModel json.Marshal(req.MCPServerIDs) failed", zap.Error(err))
	}

	mCPServerEnabledIDs, err := json.Marshal(req.MCPServerEnabledIDs)
	if err != nil {
		zap.L().Error("buildAISubAgentModel json.Marshal(req.MCPServerIDs) failed", zap.Error(err))
	}

	localTools, err := json.Marshal(req.LocalTools)
	if err != nil {
		zap.L().Error("buildAISubAgentModel json.Marshal(req.MCPServerIDs) failed", zap.Error(err))
	}
	return &gatewaymodels.AISubAgent{
		Name:                      req.Name,
		Type:                      req.Type,
		Description:               req.Description,
		Instruction:               req.Instruction,
		MCPEnabled:                req.MCPEnabled,
		MCPServerIDs:              string(mCPServerIDs),
		MCPServerEnabledIDs:       string(mCPServerEnabledIDs),
		LocalToolEnabled:          req.LocalToolEnabled,
		LocalTools:                string(localTools),
		AllowMcpServerIdsToolsRaw: req.AllowMcpServerIdsToolsRaw,
		Enable:                    req.Enable,
		DeptID:                    baizeContext.GetDeptId(c),
	}
}

func decodeSubAgentAllowMcpServerIdsTools(item *gatewaymodels.AISubAgent) {
	if item == nil {
		return
	}
	item.AllowMcpServerIdsTools = make(map[string][]string)
	content := strings.TrimSpace(item.AllowMcpServerIdsToolsRaw)
	if content == "" {
		return
	}
	_ = json.Unmarshal([]byte(content), &item.AllowMcpServerIdsTools)
}

// GetByName 查询子智能体配置详情。
func (a *AISubAgentDaoImpl) GetByName(c *gin.Context, name string) (*gatewaymodels.AISubAgent, error) {
	var item gatewaymodels.AISubAgent
	if err := a.db.WithContext(c).Table(a.table).
		Where("name = ?", name).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
