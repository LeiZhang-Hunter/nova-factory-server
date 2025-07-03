package qcRqcServiceImpl

import (
	"context"
	"fmt"
	"time"

	"nova-factory-server/app/business/qcRqc/qcRqcApi"
	"nova-factory-server/app/business/qcRqc/qcRqcDao"
	"nova-factory-server/app/business/qcRqc/qcRqcModels"
	"nova-factory-server/app/business/qcRqc/qcRqcService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QcRqcServiceImpl 退料检验单服务实现
type QcRqcServiceImpl struct {
	qcRqcDao qcRqcDao.IQcRqcDao
}

// NewQcRqcServiceImpl 创建退料检验单服务实现
func NewQcRqcServiceImpl(qcRqcDao qcRqcDao.IQcRqcDao) qcRqcService.IQcRqcService {
	return &QcRqcServiceImpl{qcRqcDao: qcRqcDao}
}

// GetQcRqcList 获取退料检验单列表
func (s *QcRqcServiceImpl) GetQcRqcList(ctx context.Context, req *qcRqcApi.QcRqcQueryReq) ([]*qcRqcApi.QcRqc, int64, error) {
	// 设置默认分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用DAO查询数据
	models, total, err := s.qcRqcDao.SelectQcRqcList(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API结构体
	apis := make([]*qcRqcApi.QcRqc, len(models))
	for i, model := range models {
		apis[i] = s.convertModelToApi(model)
	}

	return apis, total, nil
}

// GetQcRqcById 根据ID获取退料检验单
func (s *QcRqcServiceImpl) GetQcRqcById(ctx context.Context, rqcId int64) (*qcRqcApi.QcRqc, error) {
	if rqcId <= 0 {
		return nil, fmt.Errorf("退料检验单ID不能为空")
	}

	model, err := s.qcRqcDao.SelectQcRqcById(ctx, rqcId)
	if err != nil {
		return nil, err
	}

	if model == nil {
		return nil, nil
	}

	return s.convertModelToApi(model), nil
}

// CreateQcRqc 创建退料检验单
func (s *QcRqcServiceImpl) CreateQcRqc(ctx context.Context, req *qcRqcApi.QcRqcCreateReq) error {
	// 获取当前用户
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return fmt.Errorf("上下文类型错误")
	}
	userName := baizeContext.GetUserName(ginCtx)
	if userName == "" {
		return fmt.Errorf("获取用户信息失败")
	}

	// 生成检验单编号
	rqcCode := fmt.Sprintf("RQC%s%04d", time.Now().Format("20060102"), snowflake.GenID()%10000)

	// 构建退料检验单model对象
	rqc := &qcRqcModels.QcRqc{
		RqcCode:             rqcCode,
		RqcName:             req.RqcName,
		TemplateId:          req.TemplateId,
		SourceDocId:         req.SourceDocId,
		SourceDocType:       req.SourceDocType,
		SourceDocCode:       req.SourceDocCode,
		SourceLineId:        req.SourceLineId,
		ItemId:              req.ItemId,
		ItemCode:            req.ItemCode,
		ItemName:            req.ItemName,
		Specification:       req.Specification,
		UnitOfMeasure:       req.UnitOfMeasure,
		UnitName:            req.UnitName,
		BatchId:             req.BatchId,
		BatchCode:           req.BatchCode,
		QuantityCheck:       &req.QuantityCheck,
		QuantityUnqualified: &req.QuantityUnqualified,
		QuantityQualified:   req.QuantityQualified,
		CheckResult:         req.CheckResult,
		InspectDate:         req.InspectDate,
		UserId:              req.UserId,
		UserName:            req.UserName,
		NickName:            req.NickName,
		Status:              req.Status,
		Remark:              req.Remark,
		Attr1:               req.Attr1,
		Attr2:               req.Attr2,
		Attr3:               req.Attr3,
		Attr4:               req.Attr4,
		CreateBy:            &userName,
		UpdateBy:            &userName,
		CreateById:          baizeContext.GetUserId(ginCtx),
		UpdateById:          baizeContext.GetUserId(ginCtx),
	}

	// 设置默认状态
	if rqc.Status == nil {
		defaultStatus := "DRAFT"
		rqc.Status = &defaultStatus
	}

	// 设置检测日期
	if rqc.InspectDate == nil {
		now := time.Now()
		rqc.InspectDate = &now
	}

	return s.qcRqcDao.InsertQcRqc(ctx, rqc)
}

// UpdateQcRqc 更新退料检验单
func (s *QcRqcServiceImpl) UpdateQcRqc(ctx context.Context, req *qcRqcApi.QcRqcUpdateReq) error {
	// 获取当前用户
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return fmt.Errorf("上下文类型错误")
	}
	userName := baizeContext.GetUserName(ginCtx)
	if userName == "" {
		return fmt.Errorf("获取用户信息失败")
	}

	// 查询现有数据
	existingRqc, err := s.qcRqcDao.SelectQcRqcById(ctx, req.RqcId)
	if err != nil {
		zap.L().Error("查询退料检验单失败", zap.Error(err))
		return err
	}
	if existingRqc == nil {
		return fmt.Errorf("退料检验单不存在")
	}

	// 更新字段
	if req.RqcName != nil {
		existingRqc.RqcName = req.RqcName
	}
	if req.TemplateId != nil {
		existingRqc.TemplateId = *req.TemplateId
	}
	if req.SourceDocId != nil {
		existingRqc.SourceDocId = req.SourceDocId
	}
	if req.SourceDocType != nil {
		existingRqc.SourceDocType = req.SourceDocType
	}
	if req.SourceDocCode != nil {
		existingRqc.SourceDocCode = req.SourceDocCode
	}
	if req.SourceLineId != nil {
		existingRqc.SourceLineId = req.SourceLineId
	}
	if req.ItemId != nil {
		existingRqc.ItemId = *req.ItemId
	}
	if req.ItemCode != nil {
		existingRqc.ItemCode = req.ItemCode
	}
	if req.ItemName != nil {
		existingRqc.ItemName = req.ItemName
	}
	if req.Specification != nil {
		existingRqc.Specification = req.Specification
	}
	if req.UnitOfMeasure != nil {
		existingRqc.UnitOfMeasure = req.UnitOfMeasure
	}
	if req.UnitName != nil {
		existingRqc.UnitName = req.UnitName
	}
	if req.BatchId != nil {
		existingRqc.BatchId = req.BatchId
	}
	if req.BatchCode != nil {
		existingRqc.BatchCode = req.BatchCode
	}
	if req.QuantityCheck != nil {
		existingRqc.QuantityCheck = req.QuantityCheck
	}
	if req.QuantityUnqualified != nil {
		existingRqc.QuantityUnqualified = req.QuantityUnqualified
	}
	if req.QuantityQualified != nil {
		existingRqc.QuantityQualified = req.QuantityQualified
	}
	if req.CheckResult != nil {
		existingRqc.CheckResult = req.CheckResult
	}
	if req.InspectDate != nil {
		existingRqc.InspectDate = req.InspectDate
	}
	if req.UserId != nil {
		existingRqc.UserId = req.UserId
	}
	if req.UserName != nil {
		existingRqc.UserName = req.UserName
	}
	if req.NickName != nil {
		existingRqc.NickName = req.NickName
	}
	if req.Status != nil {
		existingRqc.Status = req.Status
	}
	if req.Remark != nil {
		existingRqc.Remark = req.Remark
	}
	if req.Attr1 != nil {
		existingRqc.Attr1 = req.Attr1
	}
	if req.Attr2 != nil {
		existingRqc.Attr2 = req.Attr2
	}
	if req.Attr3 != nil {
		existingRqc.Attr3 = req.Attr3
	}
	if req.Attr4 != nil {
		existingRqc.Attr4 = req.Attr4
	}

	// 更新用户信息
	existingRqc.UpdateBy = &userName
	existingRqc.UpdateById = baizeContext.GetUserId(ginCtx)

	return s.qcRqcDao.UpdateQcRqc(ctx, existingRqc)
}

// DeleteQcRqcByIds 批量删除退料检验单
func (s *QcRqcServiceImpl) DeleteQcRqcByIds(ctx context.Context, rqcIds []int64) error {
	if len(rqcIds) == 0 {
		return fmt.Errorf("退料检验单ID列表不能为空")
	}

	return s.qcRqcDao.DeleteQcRqcByIds(ctx, rqcIds)
}

// convertModelToApi 将model结构体转换为API结构体
func (s *QcRqcServiceImpl) convertModelToApi(model *qcRqcModels.QcRqc) *qcRqcApi.QcRqc {
	return &qcRqcApi.QcRqc{
		RqcId:               int64(model.ID),
		RqcCode:             model.RqcCode,
		RqcName:             model.RqcName,
		TemplateId:          model.TemplateId,
		SourceDocId:         model.SourceDocId,
		SourceDocType:       model.SourceDocType,
		SourceDocCode:       model.SourceDocCode,
		SourceLineId:        model.SourceLineId,
		ItemId:              model.ItemId,
		ItemCode:            model.ItemCode,
		ItemName:            model.ItemName,
		Specification:       model.Specification,
		UnitOfMeasure:       model.UnitOfMeasure,
		UnitName:            model.UnitName,
		BatchId:             model.BatchId,
		BatchCode:           model.BatchCode,
		QuantityCheck:       s.float64PtrToFloat64(model.QuantityCheck),
		QuantityUnqualified: s.float64PtrToFloat64(model.QuantityUnqualified),
		QuantityQualified:   model.QuantityQualified,
		CheckResult:         model.CheckResult,
		InspectDate:         model.InspectDate,
		UserId:              model.UserId,
		UserName:            model.UserName,
		NickName:            model.NickName,
		Status:              model.Status,
		Remark:              model.Remark,
		Attr1:               model.Attr1,
		Attr2:               model.Attr2,
		Attr3:               model.Attr3,
		Attr4:               model.Attr4,
		CreateBy:            s.stringPtrToString(model.CreateBy),
		CreateTime:          model.CreatedAt,
		UpdateBy:            s.stringPtrToString(model.UpdateBy),
		UpdateTime:          &model.UpdatedAt,
	}
}

// stringPtrToString 将字符串指针转换为字符串
func (s *QcRqcServiceImpl) stringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// intPtrToInt 将整数指针转换为整数
func (s *QcRqcServiceImpl) intPtrToInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// float64PtrToFloat64 将float64指针转换为float64
func (s *QcRqcServiceImpl) float64PtrToFloat64(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
