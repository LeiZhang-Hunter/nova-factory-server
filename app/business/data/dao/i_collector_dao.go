package dao

import (
	"context"
	"nova-factory-server/app/business/data/models/entity"
)

// ICollectorDAO 采集器数据访问接口
type ICollectorDAO interface {
	Create(ctx context.Context, c *entity.Collector) error
	GetByID(ctx context.Context, id string) (*entity.Collector, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*entity.Collector, error)
	List(ctx context.Context, offset, limit int, name, deviceID, status string) ([]entity.Collector, int64, error)
	ListAll(ctx context.Context) ([]entity.Collector, error)
	Update(ctx context.Context, c *entity.Collector) error
	Delete(ctx context.Context, id string) error
}
