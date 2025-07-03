package batchService

import (
	"context"
	"nova-factory-server/app/business/batch/batchApi"
)

type IBatchService interface {
	// GetBatchList 获取批次列表
	GetBatchList(ctx context.Context, req *batchApi.BatchQueryReq) ([]*batchApi.Batch, int64, error)
	// GetBatchById 根据ID获取批次
	GetBatchById(ctx context.Context, batchId int64) (*batchApi.Batch, error)
	// CreateBatch 创建批次
	CreateBatch(ctx context.Context, req *batchApi.BatchCreateReq) error
	// UpdateBatch 更新批次
	UpdateBatch(ctx context.Context, req *batchApi.BatchUpdateReq) error
	// DeleteBatch 删除批次
	DeleteBatch(ctx context.Context, req *batchApi.BatchDeleteReq) error
}
