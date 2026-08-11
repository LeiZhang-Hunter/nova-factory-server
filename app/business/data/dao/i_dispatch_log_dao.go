package dao

import (
	"context"

	"nova-factory-server/app/business/data/models/entity"

	"gorm.io/gorm"
)

// IDispatchLogDAO 下发版本记录数据访问接口
type IDispatchLogDAO interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, log *entity.DispatchLog) error
	List(ctx context.Context, offset, limit int, collectorID string) ([]entity.DispatchLog, int64, error)
	LatestByCollector(ctx context.Context, collectorID string) (*entity.DispatchLog, error)
	NextDispatchNo(ctx context.Context, collectorID string) (int, error)
}
