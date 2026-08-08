package dao

import (
	"context"

	"nova-factory-server/app/business/data/models/entity"
)

// IModelConfigDAO 模型配置数据访问接口。
type IModelConfigDAO interface {
	Get(ctx context.Context) (*entity.ModelConfig, error)
	Save(ctx context.Context, config *entity.ModelConfig) error
}
