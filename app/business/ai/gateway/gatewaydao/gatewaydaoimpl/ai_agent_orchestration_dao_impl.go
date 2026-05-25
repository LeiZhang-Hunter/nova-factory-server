package gatewaydaoimpl

import (
	"context"
	"errors"
	"nova-factory-server/app/utils/uuid"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIAgentOrchestrationDaoImpl 提供智能体编排的数据访问能力。
type AIAgentOrchestrationDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewAIAgentOrchestrationDao 创建智能体编排DAO。
func NewAIAgentOrchestrationDao(db *gorm.DB) gatewaydao.IAIAgentOrchestrationDao {
	return &AIAgentOrchestrationDaoImpl{
		db:    db,
		table: "ai_agent_orchestration",
	}
}

// Create 新增智能体编排配置。
func (a *AIAgentOrchestrationDaoImpl) Create(c *gin.Context, req *gatewaymodels.AIAgentOrchestrationUpsert) (*gatewaymodels.AIAgentOrchestration, error) {
	item, err := buildAIAgentOrchestrationModel(c, req)
	if err != nil {
		return nil, err
	}
	item.ID = snowflake.GenID()
	item.State = commonStatus.NORMAL
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// UpdateByAgentID 按智能体ID更新编排配置。
func (a *AIAgentOrchestrationDaoImpl) UpdateByAgentID(c *gin.Context, req *gatewaymodels.AIAgentOrchestrationUpsert) (*gatewaymodels.AIAgentOrchestration, error) {
	item, err := buildAIAgentOrchestrationModel(c, req)
	if err != nil {
		return nil, err
	}
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).
		Where("agent_id = ?", req.AgentID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Select("content", "config", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return a.GetByAgentID(c, req.AgentID)
}

// GetByAgentID 按智能体ID查询编排配置。
func (a *AIAgentOrchestrationDaoImpl) GetByAgentID(c *gin.Context, agentID int64) (*gatewaymodels.AIAgentOrchestration, error) {
	var item gatewaymodels.AIAgentOrchestration
	if err := a.db.WithContext(c).Table(a.table).
		Where("agent_id = ?", agentID).
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

// DeleteByAgentIDs 按智能体ID删除编排配置。
func (a *AIAgentOrchestrationDaoImpl) DeleteByAgentIDs(c *gin.Context, agentIDs []int64) error {
	now := time.Now()
	return a.db.WithContext(c).Table(a.table).
		Where("agent_id IN ?", agentIDs).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

func buildAIAgentOrchestrationModel(c *gin.Context, req *gatewaymodels.AIAgentOrchestrationUpsert) (*gatewaymodels.AIAgentOrchestration, error) {

	contentMd5 := uuid.MakeMd5([]byte(req.Content))

	configMd5 := uuid.MakeMd5([]byte(req.Config))

	return &gatewaymodels.AIAgentOrchestration{
		AgentID:    req.AgentID,
		Content:    req.Content,
		ConfigMd5:  contentMd5,
		ContentMd5: configMd5,
		Config:     req.Config,
		DeptID:     baizeContext.GetDeptId(c),
	}, nil
}

// GetConfigByAgentID 按智能体ID查询编排配置。
func (a *AIAgentOrchestrationDaoImpl) GetConfigByAgentID(c context.Context, agentID int64) (*gatewaymodels.AIAgentOrchestration, error) {
	var item gatewaymodels.AIAgentOrchestration
	if err := a.db.WithContext(c).Table(a.table).
		Where("agent_id = ?", agentID).
		Where("state = ?", commonStatus.NORMAL).
		Select("config", "config_md5").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
