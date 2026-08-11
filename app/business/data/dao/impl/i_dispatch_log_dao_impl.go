package impl

import (
	"context"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/entity"

	"gorm.io/gorm"
)

// IDispatchLogDAOImpl 下发版本记录 DAO 实现
type IDispatchLogDAOImpl struct {
	db *gorm.DB
}

// NewIDispatchLogDAOImpl 创建下发版本记录 DAO
func NewIDispatchLogDAOImpl(db *gorm.DB) dao.IDispatchLogDAO {
	return &IDispatchLogDAOImpl{db: db}
}

// CreateWithTx 在事务内创建下发记录
func (d *IDispatchLogDAOImpl) CreateWithTx(ctx context.Context, tx *gorm.DB, log *entity.DispatchLog) error {
	if tx == nil {
		tx = d.db
	}
	return tx.WithContext(ctx).Create(log).Error
}

// List 分页查询下发记录
func (d *IDispatchLogDAOImpl) List(ctx context.Context, offset, limit int, collectorID string) ([]entity.DispatchLog, int64, error) {
	var rows []entity.DispatchLog
	var total int64
	query := d.db.WithContext(ctx).Model(&entity.DispatchLog{})
	if collectorID != "" {
		query = query.Where("collector_id = ?", collectorID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// LatestByCollector 取采集器最近一次下发记录
func (d *IDispatchLogDAOImpl) LatestByCollector(ctx context.Context, collectorID string) (*entity.DispatchLog, error) {
	var log entity.DispatchLog
	if err := d.db.WithContext(ctx).Where("collector_id = ?", collectorID).Order("dispatch_no DESC").First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

// NextDispatchNo 计算采集器下一个下发序号
func (d *IDispatchLogDAOImpl) NextDispatchNo(ctx context.Context, collectorID string) (int, error) {
	var maxNo int
	err := d.db.WithContext(ctx).Raw(
		"SELECT COALESCE(MAX(dispatch_no), 0) FROM data_dispatch_logs WHERE collector_id = ?",
		collectorID,
	).Scan(&maxNo).Error
	if err != nil {
		return 0, err
	}
	return maxNo + 1, nil
}
