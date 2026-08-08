package impl

import (
	"context"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/entity"

	"gorm.io/gorm"
)

// IModelConfigDAOImpl 模型配置 DAO 实现。
type IModelConfigDAOImpl struct {
	db *gorm.DB
}

// NewIModelConfigDAOImpl 创建模型配置 DAO。
func NewIModelConfigDAOImpl(db *gorm.DB) dao.IModelConfigDAO {
	return &IModelConfigDAOImpl{db: db}
}

// Get 查询单行模型配置；不存在时返回 nil。
func (d *IModelConfigDAOImpl) Get(ctx context.Context) (*entity.ModelConfig, error) {
	var config entity.ModelConfig
	err := d.db.WithContext(ctx).
		Where("id = ?", entity.ModelConfigDefaultID).
		First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

// Save 不存在则创建，已存在则整体更新。
func (d *IModelConfigDAOImpl) Save(ctx context.Context, config *entity.ModelConfig) error {
	existing, err := d.Get(ctx)
	if err != nil {
		return err
	}
	if existing == nil {
		return d.db.WithContext(ctx).Create(config).Error
	}
	return d.db.WithContext(ctx).
		Model(&entity.ModelConfig{}).
		Where("id = ?", entity.ModelConfigDefaultID).
		Select(
			"provider", "model",
			"temperature", "enable_temperature",
			"top_p", "enable_top_p",
			"max_tokens", "enable_max_tokens",
			"max_context_count",
			"retrieval_top_k",
			"retrieval_match_threshold",
			"updated_at",
		).
		Updates(config).Error
}
