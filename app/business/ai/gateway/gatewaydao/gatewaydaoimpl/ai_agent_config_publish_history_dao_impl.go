package gatewaydaoimpl

import (
	"context"
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

// AIAgentConfigPublishHistoryDaoImpl 提供智能体配置发布历史的数据访问能力。
type AIAgentConfigPublishHistoryDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewAIAgentConfigPublishHistoryDao 创建智能体配置发布历史DAO。
func NewAIAgentConfigPublishHistoryDao(db *gorm.DB) gatewaydao.IAIAgentConfigPublishHistoryDao {
	return &AIAgentConfigPublishHistoryDaoImpl{
		db:    db,
		table: "ai_agent_config_publish_history",
	}
}

// Create 新增智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryDaoImpl) Create(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	return a.CreateWithTx(c, a.db.WithContext(c), req)
}

// CreateWithTx 在事务中新增智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryDaoImpl) CreateWithTx(c *gin.Context, tx *gorm.DB, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	item := &gatewaymodels.AIAgentConfigPublishHistory{
		ID:                 snowflake.GenID(),
		AgentID:            req.AgentID,
		Version:            req.Version,
		ConfigSnapshot:     req.ConfigSnapshot,
		ConfigMd5:          req.ConfigMd5,
		PublishDescription: req.PublishDescription,
		DeptID:             baizeContext.GetDeptId(c),
		State:              commonStatus.NORMAL,
	}
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := tx.Table(a.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Update 修改智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryDaoImpl) Update(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	return a.UpdateWithTx(c, a.db.WithContext(c), req)
}

// UpdateWithTx 在事务中修改智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryDaoImpl) UpdateWithTx(c *gin.Context, tx *gorm.DB, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	item := &gatewaymodels.AIAgentConfigPublishHistory{
		ID:                 req.ID,
		AgentID:            req.AgentID,
		ConfigMd5:          req.ConfigMd5,
		Version:            req.Version,
		ConfigSnapshot:     req.ConfigSnapshot,
		PublishDescription: req.PublishDescription,
	}
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := tx.Table(a.table).
		Where("id = ?", item.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Select("agent_id", "version", "config_snapshot", "config_md5", "publish_description", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return a.getByIDWithDB(c, tx, item.ID)
}

// GetByID 查询智能体配置发布历史详情。
func (a *AIAgentConfigPublishHistoryDaoImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	return a.getByIDWithDB(c, a.db.WithContext(c), id)
}

func (a *AIAgentConfigPublishHistoryDaoImpl) getByIDWithDB(c *gin.Context, db *gorm.DB, id int64) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	var item gatewaymodels.AIAgentConfigPublishHistory
	if err := db.Table(a.table).
		Where("id = ?", id).
		//Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByAgentIDAndVersion 按智能体ID和版本号查询发布历史。
func (a *AIAgentConfigPublishHistoryDaoImpl) GetByAgentIDAndVersion(c *gin.Context, agentID int64, version string) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	var item gatewaymodels.AIAgentConfigPublishHistory
	if err := a.db.WithContext(c).Table(a.table).
		Where("agent_id = ?", agentID).
		Where("version = ?", version).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (a *AIAgentConfigPublishHistoryDaoImpl) GetByVersion(c context.Context, version string) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	var item gatewaymodels.AIAgentConfigPublishHistory
	if err := a.db.WithContext(c).Table(a.table).
		Where("version = ?", version).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// DeleteByIDs 删除智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return a.DeleteByIDsWithTx(c, a.db.WithContext(c), ids)
}

// DeleteByIDsWithTx 在事务中删除智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryDaoImpl) DeleteByIDsWithTx(c *gin.Context, tx *gorm.DB, ids []int64) error {
	now := time.Now()
	return tx.Table(a.table).
		Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

// List 查询智能体配置发布历史列表。
func (a *AIAgentConfigPublishHistoryDaoImpl) List(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryQuery) (*gatewaymodels.AIAgentConfigPublishHistoryListData, error) {
	db := a.db.WithContext(c).Table(a.table).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req.AgentID > 0 {
		db = db.Where("agent_id = ?", req.AgentID)
	}
	if version := strings.TrimSpace(req.Version); version != "" {
		db = db.Where("version LIKE ?", "%"+version+"%")
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

	rows := make([]*gatewaymodels.AIAgentConfigPublishHistory, 0)
	if err := db.Order("create_time DESC, id DESC").
		Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &gatewaymodels.AIAgentConfigPublishHistoryListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (a *AIAgentConfigPublishHistoryDaoImpl) GetConfigsByAgentIdAndVersion(c context.Context, conditionMap map[int64]string) ([]*gatewaymodels.AIAgentConfigPublishHistory, error) {
	rows := make([]*gatewaymodels.AIAgentConfigPublishHistory, 0)
	if len(conditionMap) == 0 {
		return rows, nil
	}

	db := a.db.WithContext(c).Table(a.table)

	conditionDB := a.db.WithContext(c)
	first := true
	for agentID, version := range conditionMap {
		if first {
			conditionDB = conditionDB.Where("(agent_id = ? AND version = ?)", agentID, version)
			first = false
			continue
		}
		conditionDB = conditionDB.Or("(agent_id = ? AND version = ?)", agentID, version)
	}

	conditionDB = conditionDB.Where("state = ?", commonStatus.NORMAL)

	if err := db.Where(conditionDB).Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}
