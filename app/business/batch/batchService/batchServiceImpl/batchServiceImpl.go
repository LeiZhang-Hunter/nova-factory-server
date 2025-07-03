package batchServiceImpl

import (
	"context"
	"nova-factory-server/app/business/batch/batchApi"
	"nova-factory-server/app/business/batch/batchDao"
	"nova-factory-server/app/business/batch/batchModels"
	"nova-factory-server/app/business/batch/batchService"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BatchServiceImpl struct {
	batchDao batchDao.IBatchDao
}

func NewBatchServiceImpl(batchDao batchDao.IBatchDao) batchService.IBatchService {
	return &BatchServiceImpl{
		batchDao: batchDao,
	}
}

// GetBatchList 获取批次列表
func (s *BatchServiceImpl) GetBatchList(ctx context.Context, req *batchApi.BatchQueryReq) ([]*batchApi.Batch, int64, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	models, total, err := s.batchDao.SelectBatchList(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	apis := make([]*batchApi.Batch, len(models))
	for i, model := range models {
		apis[i] = s.convertModelToApi(model)
	}
	return apis, total, nil
}

// GetBatchById 根据ID获取批次
func (s *BatchServiceImpl) GetBatchById(ctx context.Context, batchId int64) (*batchApi.Batch, error) {
	if batchId <= 0 {
		zap.L().Error("批次ID无效", zap.Int64("batchId", batchId))
		return nil, nil
	}
	model, err := s.batchDao.SelectBatchById(ctx, batchId)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}
	return s.convertModelToApi(model), nil
}

// CreateBatch 创建批次
func (s *BatchServiceImpl) CreateBatch(ctx context.Context, req *batchApi.BatchCreateReq) error {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		zap.L().Error("无法获取gin.Context")
		return nil
	}
	batch := &batchModels.Batch{
		BatchCode:       req.BatchCode,
		ItemId:          req.ItemId,
		ItemCode:        req.ItemCode,
		ItemName:        req.ItemName,
		Specification:   req.Specification,
		UnitOfMeasure:   req.UnitOfMeasure,
		VendorId:        req.VendorId,
		VendorCode:      req.VendorCode,
		VendorName:      req.VendorName,
		VendorNick:      req.VendorNick,
		ClientId:        req.ClientId,
		ClientCode:      req.ClientCode,
		ClientName:      req.ClientName,
		ClientNick:      req.ClientNick,
		CoCode:          req.CoCode,
		PoCode:          req.PoCode,
		WorkorderId:     req.WorkorderId,
		WorkorderCode:   req.WorkorderCode,
		TaskId:          req.TaskId,
		TaskCode:        req.TaskCode,
		WorkstationId:   req.WorkstationId,
		WorkstationCode: req.WorkstationCode,
		ToolId:          req.ToolId,
		ToolCode:        req.ToolCode,
		MoldId:          req.MoldId,
		MoldCode:        req.MoldCode,
		LotNumber:       req.LotNumber,
		QualityStatus:   req.QualityStatus,
		Remark:          req.Remark,
		Attr1:           req.Attr1,
		Attr2:           req.Attr2,
		Attr3:           req.Attr3,
		Attr4:           req.Attr4,
		CreateBy:        s.stringToPtr(baizeContext.GetUserName(ginCtx)),
		UpdateBy:        s.stringToPtr(baizeContext.GetUserName(ginCtx)),
		CreateById:      s.int64ToPtr(baizeContext.GetUserId(ginCtx)),
		UpdateById:      s.int64ToPtr(baizeContext.GetUserId(ginCtx)),
	}
	return s.batchDao.InsertBatch(ctx, batch)
}

// UpdateBatch 更新批次
func (s *BatchServiceImpl) UpdateBatch(ctx context.Context, req *batchApi.BatchUpdateReq) error {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		zap.L().Error("无法获取gin.Context")
		return nil
	}
	original, err := s.batchDao.SelectBatchById(ctx, req.BatchId)
	if err != nil {
		zap.L().Error("查询原批次失败", zap.Int64("batchId", req.BatchId), zap.Error(err))
		return err
	}
	if original == nil {
		zap.L().Error("批次不存在", zap.Int64("batchId", req.BatchId))
		return nil
	}
	batch := &batchModels.Batch{
		Model:           original.Model,
		BatchCode:       req.BatchCode,
		ItemId:          req.ItemId,
		ItemCode:        req.ItemCode,
		ItemName:        req.ItemName,
		Specification:   req.Specification,
		UnitOfMeasure:   req.UnitOfMeasure,
		VendorId:        req.VendorId,
		VendorCode:      req.VendorCode,
		VendorName:      req.VendorName,
		VendorNick:      req.VendorNick,
		ClientId:        req.ClientId,
		ClientCode:      req.ClientCode,
		ClientName:      req.ClientName,
		ClientNick:      req.ClientNick,
		CoCode:          req.CoCode,
		PoCode:          req.PoCode,
		WorkorderId:     req.WorkorderId,
		WorkorderCode:   req.WorkorderCode,
		TaskId:          req.TaskId,
		TaskCode:        req.TaskCode,
		WorkstationId:   req.WorkstationId,
		WorkstationCode: req.WorkstationCode,
		ToolId:          req.ToolId,
		ToolCode:        req.ToolCode,
		MoldId:          req.MoldId,
		MoldCode:        req.MoldCode,
		LotNumber:       req.LotNumber,
		QualityStatus:   req.QualityStatus,
		Remark:          req.Remark,
		Attr1:           req.Attr1,
		Attr2:           req.Attr2,
		Attr3:           req.Attr3,
		Attr4:           req.Attr4,
		CreateBy:        original.CreateBy,
		CreateById:      original.CreateById,
		UpdateBy:        s.stringToPtr(baizeContext.GetUserName(ginCtx)),
		UpdateById:      s.int64ToPtr(baizeContext.GetUserId(ginCtx)),
	}
	return s.batchDao.UpdateBatch(ctx, batch)
}

// DeleteBatch 删除批次
func (s *BatchServiceImpl) DeleteBatch(ctx context.Context, req *batchApi.BatchDeleteReq) error {
	if len(req.BatchIds) == 0 {
		zap.L().Error("批次ID列表为空")
		return nil
	}
	return s.batchDao.DeleteBatchByIds(ctx, req.BatchIds)
}

// convertModelToApi 将model结构体转换为API结构体
func (s *BatchServiceImpl) convertModelToApi(model *batchModels.Batch) *batchApi.Batch {
	return &batchApi.Batch{
		BatchId:         int64(model.ID),
		BatchCode:       model.BatchCode,
		ItemId:          model.ItemId,
		ItemCode:        model.ItemCode,
		ItemName:        model.ItemName,
		Specification:   model.Specification,
		UnitOfMeasure:   model.UnitOfMeasure,
		VendorId:        model.VendorId,
		VendorCode:      model.VendorCode,
		VendorName:      model.VendorName,
		VendorNick:      model.VendorNick,
		ClientId:        model.ClientId,
		ClientCode:      model.ClientCode,
		ClientName:      model.ClientName,
		ClientNick:      model.ClientNick,
		CoCode:          model.CoCode,
		PoCode:          model.PoCode,
		WorkorderId:     model.WorkorderId,
		WorkorderCode:   model.WorkorderCode,
		TaskId:          model.TaskId,
		TaskCode:        model.TaskCode,
		WorkstationId:   model.WorkstationId,
		WorkstationCode: model.WorkstationCode,
		ToolId:          model.ToolId,
		ToolCode:        model.ToolCode,
		MoldId:          model.MoldId,
		MoldCode:        model.MoldCode,
		LotNumber:       model.LotNumber,
		QualityStatus:   model.QualityStatus,
		Remark:          model.Remark,
		Attr1:           model.Attr1,
		Attr2:           model.Attr2,
		Attr3:           model.Attr3,
		Attr4:           model.Attr4,
		CreateBy:        model.CreateBy,
		UpdateBy:        model.UpdateBy,
		CreateTime:      &model.CreatedAt,
		UpdateTime:      &model.UpdatedAt,
	}
}

func (s *BatchServiceImpl) stringToPtr(str string) *string {
	return &str
}

func (s *BatchServiceImpl) int64ToPtr(val int64) *int64 {
	return &val
}
