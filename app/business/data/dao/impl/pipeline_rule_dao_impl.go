package impl

import (
	"context"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/entity"

	"gorm.io/gorm"
)

type PipelineRuleDAOImpl struct{ db *gorm.DB }
type ServiceConnectionDAOImpl struct{ db *gorm.DB }

func NewPipelineRuleDAO(db *gorm.DB) dao.IPipelineRuleDAO { return &PipelineRuleDAOImpl{db: db} }
func NewServiceConnectionDAO(db *gorm.DB) dao.IServiceConnectionDAO {
	return &ServiceConnectionDAOImpl{db: db}
}

func (d *PipelineRuleDAOImpl) Create(ctx context.Context, rule *entity.PipelineRule) error {
	return d.db.WithContext(ctx).Create(rule).Error
}
func (d *PipelineRuleDAOImpl) GetByID(ctx context.Context, id string) (*entity.PipelineRule, error) {
	var rule entity.PipelineRule
	if err := d.db.WithContext(ctx).First(&rule, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}
func (d *PipelineRuleDAOImpl) GetByName(ctx context.Context, sourceType, name string) (*entity.PipelineRule, error) {
	var rule entity.PipelineRule
	if err := d.db.WithContext(ctx).Where("source_type = ? AND name = ?", sourceType, name).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}
func (d *PipelineRuleDAOImpl) List(ctx context.Context, offset, limit int, sourceType, name, status string) ([]entity.PipelineRule, int64, error) {
	var rows []entity.PipelineRule
	var total int64
	query := d.db.WithContext(ctx).Model(&entity.PipelineRule{})
	if sourceType != "" {
		query = query.Where("source_type = ?", sourceType)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}
func (d *PipelineRuleDAOImpl) Update(ctx context.Context, rule *entity.PipelineRule) error {
	return d.db.WithContext(ctx).Model(rule).Select("name", "description", "status", "version", "config", "updated_at").Updates(rule).Error
}
func (d *PipelineRuleDAOImpl) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&entity.PipelineRule{}, "id = ?", id).Error
}
func (d *PipelineRuleDAOImpl) CountConnectionReferences(ctx context.Context, id string) (int64, error) {
	var total int64
	err := d.db.WithContext(ctx).Model(&entity.PipelineRule{}).
		Where("JSON_SEARCH(config, 'one', ?, NULL, '$.pipelines[*].source.connectionRef') IS NOT NULL", id).
		Count(&total).Error
	return total, err
}

func (d *ServiceConnectionDAOImpl) Create(ctx context.Context, row *entity.ServiceConnection) error {
	return d.db.WithContext(ctx).Create(row).Error
}
func (d *ServiceConnectionDAOImpl) GetByID(ctx context.Context, id string) (*entity.ServiceConnection, error) {
	var row entity.ServiceConnection
	if err := d.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
func (d *ServiceConnectionDAOImpl) GetByName(ctx context.Context, sourceType, name string) (*entity.ServiceConnection, error) {
	var row entity.ServiceConnection
	if err := d.db.WithContext(ctx).Where("source_type = ? AND name = ?", sourceType, name).First(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
func (d *ServiceConnectionDAOImpl) List(ctx context.Context, offset, limit int, sourceType, name, status string) ([]entity.ServiceConnection, int64, error) {
	var rows []entity.ServiceConnection
	var total int64
	query := d.db.WithContext(ctx).Model(&entity.ServiceConnection{})
	if sourceType != "" {
		query = query.Where("source_type = ?", sourceType)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}
func (d *ServiceConnectionDAOImpl) Update(ctx context.Context, row *entity.ServiceConnection) error {
	return d.db.WithContext(ctx).Model(row).Select("name", "description", "status", "config", "encrypted_credentials", "updated_at").Updates(row).Error
}
func (d *ServiceConnectionDAOImpl) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&entity.ServiceConnection{}, "id = ?", id).Error
}
