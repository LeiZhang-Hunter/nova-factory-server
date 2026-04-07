package gatewaydaoimpl

import (
	"errors"
	"nova-factory-server/app/utils/snowflake"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AIGatewayDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAIGatewayDao(db *gorm.DB) gatewaydao.IAIGatewayDao {
	return &AIGatewayDaoImpl{
		db:    db,
		table: "ai_gateway",
	}
}

func (a *AIGatewayDaoImpl) Create(c *gin.Context, req *gatewaymodels.AIGatewayUpsert) (*gatewaymodels.AIGateway, error) {
	item := &gatewaymodels.AIGateway{
		ID:      snowflake.GenID(),
		Name:    req.Name,
		BaseURL: req.BaseURL,
		APIKey:  req.APIKey,
		Enabled: req.Enabled,
		DeptID:  baizeContext.GetDeptId(c),
		State:   commonStatus.NORMAL,
	}
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (a *AIGatewayDaoImpl) Update(c *gin.Context, req *gatewaymodels.AIGatewayUpsert) (*gatewaymodels.AIGateway, error) {
	item := &gatewaymodels.AIGateway{
		ID:      req.ID,
		Name:    req.Name,
		BaseURL: req.BaseURL,
		APIKey:  req.APIKey,
		Enabled: req.Enabled,
	}
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", item.ID).Where("state = ?", commonStatus.NORMAL).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Select("name", "base_url", "api_key", "enabled", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return a.GetByID(c, item.ID)
}

func (a *AIGatewayDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return a.db.WithContext(c).Table(a.table).Where("id IN ?", ids).Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

func (a *AIGatewayDaoImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AIGateway, error) {
	var item gatewaymodels.AIGateway
	if err := a.db.WithContext(c).Table(a.table).Where("id = ?", id).Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (a *AIGatewayDaoImpl) List(c *gin.Context, req *gatewaymodels.AIGatewayQuery) (*gatewaymodels.AIGatewayListData, error) {
	db := a.db.Table(a.table)
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Enabled != nil {
		db = db.Where("enabled = ?", req.Enabled)
	}
	if req.Active != nil {
		db = db.Where("active = ?", req.Active)
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
	rows := make([]*gatewaymodels.AIGateway, 0)
	if err := db.Order("id DESC").Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &gatewaymodels.AIGatewayListData{
		Rows:  rows,
		Total: total,
	}, nil
}
