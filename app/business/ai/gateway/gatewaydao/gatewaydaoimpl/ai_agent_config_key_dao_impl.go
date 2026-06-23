package gatewaydaoimpl

import (
	"errors"
	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AgentConfigKeyDaoImpl API Key DAO 实现。
type AgentConfigKeyDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewAgentConfigKeyDao 创建 AgentConfigKeyDaoImpl。
func NewAgentConfigKeyDao(db *gorm.DB) gatewaydao.IAgentConfigKeyDao {
	return &AgentConfigKeyDaoImpl{
		db:    db,
		table: "ai_agent_config_key",
	}
}

func (a *AgentConfigKeyDaoImpl) Create(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	item := &gatewaymodels.AgentConfigKey{
		ID:     snowflake.GenID(),
		Key:    req.Key,
		DeptID: baizeContext.GetDeptId(c),
		State:  commonStatus.NORMAL,
	}
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (a *AgentConfigKeyDaoImpl) Update(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	item := &gatewaymodels.AgentConfigKey{
		ID:  req.ID,
		Key: req.Key,
	}
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", item.ID).Where("state = ?", commonStatus.NORMAL).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Select("key", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return a.GetByID(c, item.ID)
}

func (a *AgentConfigKeyDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	var key gatewaymodels.AgentConfigKey
	return a.db.WithContext(c).Table(a.table).Where("id IN ?", ids).Where("create_by = ?", baizeContext.GetUserId(c)).
		Delete(&key).Error
}

func (a *AgentConfigKeyDaoImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AgentConfigKey, error) {
	var item gatewaymodels.AgentConfigKey
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (a *AgentConfigKeyDaoImpl) GetByKey(c *gin.Context, key string) (*gatewaymodels.AgentConfigKey, error) {
	var item gatewaymodels.AgentConfigKey
	if err := a.db.WithContext(c).Table(a.table).Where("`key` = ?", key).
		Where("state = ?", commonStatus.NORMAL).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (a *AgentConfigKeyDaoImpl) List(c *gin.Context, req *gatewaymodels.AgentConfigKeyQuery) (*gatewaymodels.AgentConfigKeyListData, error) {
	db := a.db.Table(a.table)
	if req.Key != "" {
		db = db.Where("`key` LIKE ?", "%"+req.Key+"%")
	}
	db = db.Where("state = ?", commonStatus.NORMAL)
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
	rows := make([]*gatewaymodels.AgentConfigKey, 0)
	if err := db.Order("id DESC").Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &gatewaymodels.AgentConfigKeyListData{
		Rows:  rows,
		Total: total,
	}, nil
}
