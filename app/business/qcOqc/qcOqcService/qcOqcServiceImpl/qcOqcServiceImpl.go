package qcOqcServiceImpl

import (
	"context"

	"nova-factory-server/app/business/qcOqc/qcOqcApi"
	"nova-factory-server/app/business/qcOqc/qcOqcDao"
	"nova-factory-server/app/business/qcOqc/qcOqcModels"
	"nova-factory-server/app/business/qcOqc/qcOqcService"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QcOqcServiceImpl 出货检验单服务实现
type QcOqcServiceImpl struct {
	qcOqcDao qcOqcDao.IQcOqcDao
}

// NewQcOqcServiceImpl 创建出货检验单服务实现
func NewQcOqcServiceImpl(qcOqcDao qcOqcDao.IQcOqcDao) qcOqcService.IQcOqcService {
	return &QcOqcServiceImpl{
		qcOqcDao: qcOqcDao,
	}
}

// GetQcOqcList 获取出货检验单列表
func (s *QcOqcServiceImpl) GetQcOqcList(ctx context.Context, req *qcOqcApi.QcOqcQueryReq) ([]*qcOqcApi.QcOqc, int64, error) {
	// 设置默认分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 调用DAO查询数据
	models, total, err := s.qcOqcDao.SelectQcOqcList(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API结构体
	apis := make([]*qcOqcApi.QcOqc, len(models))
	for i, model := range models {
		apis[i] = s.convertModelToApi(model)
	}

	return apis, total, nil
}

// GetQcOqcById 根据ID获取出货检验单
func (s *QcOqcServiceImpl) GetQcOqcById(ctx context.Context, oqcId int64) (*qcOqcApi.QcOqc, error) {
	if oqcId <= 0 {
		zap.L().Error("出货检验单ID无效", zap.Int64("oqcId", oqcId))
		return nil, nil
	}

	model, err := s.qcOqcDao.SelectQcOqcById(ctx, oqcId)
	if err != nil {
		return nil, err
	}

	if model == nil {
		return nil, nil
	}

	return s.convertModelToApi(model), nil
}

// CreateQcOqc 创建出货检验单
func (s *QcOqcServiceImpl) CreateQcOqc(ctx context.Context, req *qcOqcApi.QcOqcCreateReq) error {
	// 从context中获取gin.Context
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		zap.L().Error("无法获取gin.Context")
		return nil
	}

	// 构建出货检验单model对象
	oqc := &qcOqcModels.QcOqc{
		OqcCode:                req.OqcCode,
		OqcName:                &req.OqcName,
		TemplateId:             req.TemplateId,
		SourceDocId:            req.SourceDocId,
		SourceDocType:          req.SourceDocType,
		SourceDocCode:          req.SourceDocCode,
		SourceLineId:           req.SourceLineId,
		ClientId:               req.ClientId,
		ClientCode:             req.ClientCode,
		ClientName:             req.ClientName,
		BatchCode:              req.BatchCode,
		ItemId:                 req.ItemId,
		ItemCode:               req.ItemCode,
		ItemName:               req.ItemName,
		Specification:          req.Specification,
		UnitOfMeasure:          req.UnitOfMeasure,
		QuantityMinCheck:       req.QuantityMinCheck,
		QuantityMaxUnqualified: req.QuantityMaxUnqualified,
		QuantityOut:            req.QuantityOut,
		QuantityCheck:          req.QuantityCheck,
		QuantityUnqualified:    req.QuantityUnqualified,
		QuantityQualified:      req.QuantityQualified,
		CrRate:                 req.CrRate,
		MajRate:                req.MajRate,
		MinRate:                req.MinRate,
		CrQuantity:             req.CrQuantity,
		MajQuantity:            req.MajQuantity,
		MinQuantity:            req.MinQuantity,
		CheckResult:            req.CheckResult,
		OutDate:                req.OutDate,
		InspectDate:            req.InspectDate,
		Inspector:              req.Inspector,
		Status:                 req.Status,
		Remark:                 &req.Remark,
		Attr1:                  req.Attr1,
		Attr2:                  req.Attr2,
		Attr3:                  req.Attr3,
		Attr4:                  req.Attr4,
		CreateBy:               s.stringToPtr(baizeContext.GetUserName(ginCtx)),
		UpdateBy:               s.stringToPtr(baizeContext.GetUserName(ginCtx)),
		CreateById:             baizeContext.GetUserId(ginCtx),
		UpdateById:             baizeContext.GetUserId(ginCtx),
	}

	// 设置默认值
	if oqc.Status == nil {
		defaultStatus := "PREPARE"
		oqc.Status = &defaultStatus
	}

	// 调用DAO创建数据
	return s.qcOqcDao.InsertQcOqc(ctx, oqc)
}

// UpdateQcOqc 更新出货检验单
func (s *QcOqcServiceImpl) UpdateQcOqc(ctx context.Context, req *qcOqcApi.QcOqcUpdateReq) error {
	// 从context中获取gin.Context
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		zap.L().Error("无法获取gin.Context")
		return nil
	}

	// 先查询原记录
	original, err := s.qcOqcDao.SelectQcOqcById(ctx, req.OqcId)
	if err != nil {
		zap.L().Error("查询原出货检验单失败", zap.Int64("oqcId", req.OqcId), zap.Error(err))
		return err
	}
	if original == nil {
		zap.L().Error("出货检验单不存在", zap.Int64("oqcId", req.OqcId))
		return nil
	}

	// 构建更新对象
	oqc := &qcOqcModels.QcOqc{
		Model:                  original.Model,
		OqcCode:                req.OqcCode,
		OqcName:                &req.OqcName,
		TemplateId:             req.TemplateId,
		SourceDocId:            req.SourceDocId,
		SourceDocType:          req.SourceDocType,
		SourceDocCode:          req.SourceDocCode,
		SourceLineId:           req.SourceLineId,
		ClientId:               req.ClientId,
		ClientCode:             req.ClientCode,
		ClientName:             req.ClientName,
		BatchCode:              req.BatchCode,
		ItemId:                 req.ItemId,
		ItemCode:               req.ItemCode,
		ItemName:               req.ItemName,
		Specification:          req.Specification,
		UnitOfMeasure:          req.UnitOfMeasure,
		QuantityMinCheck:       req.QuantityMinCheck,
		QuantityMaxUnqualified: req.QuantityMaxUnqualified,
		QuantityOut:            req.QuantityOut,
		QuantityCheck:          req.QuantityCheck,
		QuantityUnqualified:    req.QuantityUnqualified,
		QuantityQualified:      req.QuantityQualified,
		CrRate:                 req.CrRate,
		MajRate:                req.MajRate,
		MinRate:                req.MinRate,
		CrQuantity:             req.CrQuantity,
		MajQuantity:            req.MajQuantity,
		MinQuantity:            req.MinQuantity,
		CheckResult:            req.CheckResult,
		OutDate:                req.OutDate,
		InspectDate:            req.InspectDate,
		Inspector:              req.Inspector,
		Status:                 req.Status,
		Remark:                 &req.Remark,
		Attr1:                  req.Attr1,
		Attr2:                  req.Attr2,
		Attr3:                  req.Attr3,
		Attr4:                  req.Attr4,
		CreateBy:               original.CreateBy,
		CreateById:             original.CreateById,
		UpdateBy:               s.stringToPtr(baizeContext.GetUserName(ginCtx)),
		UpdateById:             baizeContext.GetUserId(ginCtx),
	}

	// 调用DAO更新数据
	return s.qcOqcDao.UpdateQcOqc(ctx, oqc)
}

// DeleteQcOqc 删除出货检验单
func (s *QcOqcServiceImpl) DeleteQcOqc(ctx context.Context, req *qcOqcApi.QcOqcDeleteReq) error {
	if len(req.OqcIds) == 0 {
		zap.L().Error("出货检验单ID列表为空")
		return nil
	}

	// 调用DAO删除数据
	return s.qcOqcDao.DeleteQcOqcByIds(ctx, req.OqcIds)
}

// convertModelToApi 将model结构体转换为API结构体
func (s *QcOqcServiceImpl) convertModelToApi(model *qcOqcModels.QcOqc) *qcOqcApi.QcOqc {
	return &qcOqcApi.QcOqc{
		OqcId:                  int64(model.ID),
		OqcCode:                model.OqcCode,
		OqcName:                s.stringPtrToString(model.OqcName),
		TemplateId:             model.TemplateId,
		SourceDocId:            model.SourceDocId,
		SourceDocType:          model.SourceDocType,
		SourceDocCode:          model.SourceDocCode,
		SourceLineId:           model.SourceLineId,
		ClientId:               model.ClientId,
		ClientCode:             model.ClientCode,
		ClientName:             model.ClientName,
		BatchCode:              model.BatchCode,
		ItemId:                 model.ItemId,
		ItemCode:               model.ItemCode,
		ItemName:               model.ItemName,
		Specification:          model.Specification,
		UnitOfMeasure:          model.UnitOfMeasure,
		QuantityMinCheck:       model.QuantityMinCheck,
		QuantityMaxUnqualified: model.QuantityMaxUnqualified,
		QuantityOut:            model.QuantityOut,
		QuantityCheck:          model.QuantityCheck,
		QuantityUnqualified:    model.QuantityUnqualified,
		QuantityQualified:      model.QuantityQualified,
		CrRate:                 model.CrRate,
		MajRate:                model.MajRate,
		MinRate:                model.MinRate,
		CrQuantity:             model.CrQuantity,
		MajQuantity:            model.MajQuantity,
		MinQuantity:            model.MinQuantity,
		CheckResult:            model.CheckResult,
		OutDate:                model.OutDate,
		InspectDate:            model.InspectDate,
		Inspector:              model.Inspector,
		Status:                 model.Status,
		Remark:                 s.stringPtrToString(model.Remark),
		Attr1:                  model.Attr1,
		Attr2:                  model.Attr2,
		Attr3:                  model.Attr3,
		Attr4:                  model.Attr4,
		CreateBy:               s.stringPtrToString(model.CreateBy),
		CreateTime:             &model.CreatedAt,
		UpdateBy:               s.stringPtrToString(model.UpdateBy),
		UpdateTime:             &model.UpdatedAt,
	}
}

// stringPtrToString 将字符串指针转换为字符串
func (s *QcOqcServiceImpl) stringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// stringToPtr 将字符串转换为字符串指针
func (s *QcOqcServiceImpl) stringToPtr(str string) *string {
	return &str
}
