package impl

import (
	"context"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/entity"

	"gorm.io/gorm"
)

// ICollectorDAOImpl 采集器 DAO 实现
type ICollectorDAOImpl struct {
	db *gorm.DB
}

// NewICollectorDAOImpl 创建采集器 DAO
func NewICollectorDAOImpl(db *gorm.DB) dao.ICollectorDAO {
	return &ICollectorDAOImpl{db: db}
}

// Create 创建采集器
func (d *ICollectorDAOImpl) Create(ctx context.Context, c *entity.Collector) error {
	return d.db.WithContext(ctx).Create(c).Error
}

// GetByID 按 ID 查询
func (d *ICollectorDAOImpl) GetByID(ctx context.Context, id string) (*entity.Collector, error) {
	var c entity.Collector
	err := d.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// GetByDeviceID 按设备 ID 查询
func (d *ICollectorDAOImpl) GetByDeviceID(ctx context.Context, deviceID string) (*entity.Collector, error) {
	var c entity.Collector
	err := d.db.WithContext(ctx).
		Where("device_id = ? AND deleted_at IS NULL", deviceID).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// List 分页列表
func (d *ICollectorDAOImpl) List(ctx context.Context, offset, limit int, name, deviceID, status string) ([]entity.Collector, int64, error) {
	var collectors []entity.Collector
	var total int64

	query := d.db.WithContext(ctx).
		Model(&entity.Collector{}).
		Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if deviceID != "" {
		query = query.Where("device_id LIKE ?", "%"+deviceID+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&collectors).Error

	return collectors, total, err
}

// ListAll 获取全部采集器
func (d *ICollectorDAOImpl) ListAll(ctx context.Context) ([]entity.Collector, error) {
	var collectors []entity.Collector
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&collectors).Error
	return collectors, err
}

// Update 更新采集器
func (d *ICollectorDAOImpl) Update(ctx context.Context, c *entity.Collector) error {
	return d.db.WithContext(ctx).
		Model(&entity.Collector{}).
		Where("id = ? AND deleted_at IS NULL", c.ID).
		Updates(c).Error
}

// Delete 软删除
func (d *ICollectorDAOImpl) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).
		Model(&entity.Collector{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
