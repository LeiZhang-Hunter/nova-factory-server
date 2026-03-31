package settingdaoimpl

import (
	"errors"
	"time"

	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AgentConfigDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAgentConfigDao(db *gorm.DB) settingdao.IAgentConfigDao {
	return &AgentConfigDaoImpl{
		db:    db,
		table: "erp_sales_order_agent_config",
	}
}

func (a *AgentConfigDaoImpl) Create(c *gin.Context, req *settingmodels.AgentConfigUpsert) (*settingmodels.AgentConfig, error) {
	did := baizeContext.GetDeptId(c)
	model := &settingmodels.AgentConfig{
		Name:    req.Name,
		AgentID: req.AgentID,
		Remark:  req.Remark,
		Status:  req.Status,
		DeptID:  did,
		State:   0,
	}
	model.ID = 1
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (a *AgentConfigDaoImpl) Update(c *gin.Context, req *settingmodels.AgentConfigUpsert) (*settingmodels.AgentConfig, error) {
	model := &settingmodels.AgentConfig{
		ID:      1,
		Name:    req.Name,
		AgentID: req.AgentID,
		Remark:  req.Remark,
		Status:  req.Status,
	}
	model.SetUpdateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", model.ID).Where("state = 0").
		Select("name", "agent_id", "remark", "status", "update_by", "update_time").Updates(model).Error; err != nil {
		return nil, err
	}
	return a.GetByID(c, int64(req.ID))
}

func (a *AgentConfigDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	uid := baizeContext.GetUserId(c)
	now := time.Now()
	return a.db.WithContext(c).Table(a.table).Where("id IN ?", ids).Updates(map[string]interface{}{
		"state":       -1,
		"update_by":   uid,
		"update_time": now,
	}).Error
}

func (a *AgentConfigDaoImpl) GetByID(c *gin.Context, id int64) (*settingmodels.AgentConfig, error) {
	var item settingmodels.AgentConfig
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", id).Where("state = 0").First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return &item, nil
}

func (a *AgentConfigDaoImpl) List(c *gin.Context, req *settingmodels.AgentConfigQuery) (*settingmodels.AgentConfigListData, error) {
	db := a.db.WithContext(c).Table(a.table).Where("state = 0")
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.AgentID != "" {
		db = db.Where("agent_id = ?", req.AgentID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", req.Status)
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
	rows := make([]*settingmodels.AgentConfig, 0)
	if err := db.Order("id DESC").Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &settingmodels.AgentConfigListData{
		Rows:  rows,
		Total: total,
	}, nil
}
