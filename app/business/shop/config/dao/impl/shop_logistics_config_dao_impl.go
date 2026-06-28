package impl

import (
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShopLogisticsConfigDaoImpl 物流配置数据访问实现
type ShopLogisticsConfigDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewShopLogisticsConfigDao 创建物流配置 DAO
func NewShopLogisticsConfigDao(db *gorm.DB) dao.IShopLogisticsConfigDao {
	return &ShopLogisticsConfigDaoImpl{
		db:    db,
		table: "shop_logistics_config",
	}
}

// Create 新增物流配置记录
func (d *ShopLogisticsConfigDaoImpl) Create(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error) {
	model := &models.ShopLogisticsConfig{
		ID:     snowflake.GenID(),
		Type:   strings.TrimSpace(req.Type),
		Data:   req.Data,
		Status: req.Status,
		DeptID: baizeContext.GetDeptId(c),
		State:  commonStatus.NORMAL,
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table(d.table).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

// Update 修改物流配置记录
func (d *ShopLogisticsConfigDaoImpl) Update(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error) {
	model := &models.ShopLogisticsConfig{
		ID:     req.ID,
		Type:   strings.TrimSpace(req.Type),
		Data:   req.Data,
		Status: req.Status,
	}
	model.SetUpdateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table(d.table).
		Where("id = ?", req.ID).
		Where("state = 0").
		Select("type", "data", "status", "update_by", "update_time").
		Updates(model).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, int64(req.ID))
}

// DeleteByIDs 软删除物流配置记录
func (d *ShopLogisticsConfigDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return d.db.WithContext(c).Table(d.table).Where("id IN ?", ids).Updates(map[string]interface{}{
		"state":       -1,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": now,
	}).Error
}

// GetByID 根据主键查询物流配置
func (d *ShopLogisticsConfigDaoImpl) GetByID(c *gin.Context, id int64) (*models.ShopLogisticsConfig, error) {
	var item models.ShopLogisticsConfig
	if err := d.db.WithContext(c).Table(d.table).
		Where("id = ?", id).
		Where("state = 0").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByType 根据类型查询物流配置
func (d *ShopLogisticsConfigDaoImpl) GetByType(c *gin.Context, typ string) (*models.ShopLogisticsConfig, error) {
	var item models.ShopLogisticsConfig
	if err := d.db.WithContext(c).Table(d.table).
		Where("type = ?", strings.TrimSpace(typ)).
		Where("state = 0").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// List 分页查询物流配置
func (d *ShopLogisticsConfigDaoImpl) List(c *gin.Context, req *models.ShopLogisticsConfigQuery) (*models.ShopLogisticsConfigListData, error) {
	q := d.db.WithContext(c).Table(d.table).Where("state = 0")
	if req.Type != "" {
		q = q.Where("type LIKE ?", "%"+strings.TrimSpace(req.Type)+"%")
	}
	if req.Status != nil {
		q = q.Where("status = ?", req.Status)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.ShopLogisticsConfig, 0)
	if err := q.Order("id DESC").
		Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.ShopLogisticsConfigListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (i *ShopLogisticsConfigDaoImpl) GetEnabled(c *gin.Context) (*models.ShopLogisticsConfig, error) {
	var item models.ShopLogisticsConfig
	err := i.db.WithContext(c).Table(i.table).
		Where("status = ?", true).
		Where("state = ?", commonStatus.NORMAL).
		Order("id DESC").
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}
