package batchDaoImpl

import (
	"context"
	"nova-factory-server/app/business/batch/batchApi"
	"nova-factory-server/app/business/batch/batchDao"
	"nova-factory-server/app/business/batch/batchModels"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BatchDaoImpl struct {
	db *gorm.DB
}

func NewBatchDaoImpl(db *gorm.DB) batchDao.IBatchDao {
	return &BatchDaoImpl{db: db}
}

func (dao *BatchDaoImpl) SelectBatchList(ctx context.Context, req *batchApi.BatchQueryReq) ([]*batchModels.Batch, int64, error) {
	var result []*batchModels.Batch
	var total int64
	query := dao.db.WithContext(ctx).Model(&batchModels.Batch{})

	// 应用查询条件
	if req.BatchCode != nil && *req.BatchCode != "" {
		query = query.Where("batch_code LIKE ?", "%"+*req.BatchCode+"%")
	}
	if req.ItemId != nil {
		query = query.Where("item_id = ?", *req.ItemId)
	}
	if req.ItemCode != nil && *req.ItemCode != "" {
		query = query.Where("item_code LIKE ?", "%"+*req.ItemCode+"%")
	}
	if req.ItemName != nil && *req.ItemName != "" {
		query = query.Where("item_name LIKE ?", "%"+*req.ItemName+"%")
	}
	if req.Specification != nil && *req.Specification != "" {
		query = query.Where("specification LIKE ?", "%"+*req.Specification+"%")
	}

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("查询批次总数失败", zap.Error(err))
		return nil, 0, err
	}

	// 查询数据
	err = query.Order("created_at DESC").
		Limit(req.PageSize).
		Offset((req.PageNum - 1) * req.PageSize).
		Find(&result).Error
	if err != nil {
		zap.L().Error("查询批次列表失败", zap.Error(err))
		return nil, 0, err
	}
	return result, total, nil
}

func (dao *BatchDaoImpl) SelectBatchById(ctx context.Context, id int64) (*batchModels.Batch, error) {
	var batch batchModels.Batch
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&batch).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		zap.L().Error("根据ID查询批次失败", zap.Error(err))
		return nil, err
	}
	return &batch, nil
}

func (dao *BatchDaoImpl) InsertBatch(ctx context.Context, batch *batchModels.Batch) error {
	return dao.db.WithContext(ctx).Create(batch).Error
}

func (dao *BatchDaoImpl) UpdateBatch(ctx context.Context, batch *batchModels.Batch) error {
	return dao.db.WithContext(ctx).Save(batch).Error
}

func (dao *BatchDaoImpl) DeleteBatchByIds(ctx context.Context, ids []int64) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&batchModels.Batch{}).Error
}
