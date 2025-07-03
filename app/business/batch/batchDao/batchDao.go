package batchDao

import (
	"context"
	"nova-factory-server/app/business/batch/batchApi"
	"nova-factory-server/app/business/batch/batchModels"
)

type IBatchDao interface {
	// SelectBatchList 查询批次列表
	SelectBatchList(ctx context.Context, req *batchApi.BatchQueryReq) ([]*batchModels.Batch, int64, error)
	// SelectBatchById 根据ID查询批次
	SelectBatchById(ctx context.Context, id int64) (*batchModels.Batch, error)
	// InsertBatch 新增批次
	InsertBatch(ctx context.Context, batch *batchModels.Batch) error
	// UpdateBatch 更新批次
	UpdateBatch(ctx context.Context, batch *batchModels.Batch) error
	// DeleteBatchByIds 批量删除批次
	DeleteBatchByIds(ctx context.Context, ids []int64) error
}
