package qcIpqcServiceImpl

import (
	"context"
	"fmt"
	"nova-factory-server/app/business/qcIpqc/qcIpqcApi"
	"nova-factory-server/app/business/qcIpqc/qcIpqcDao"
	"nova-factory-server/app/business/qcIpqc/qcIpqcModels"
	"nova-factory-server/app/business/qcIpqc/qcIpqcService"

	"github.com/gin-gonic/gin"
)

type QcIpqcServiceImpl struct {
	dao qcIpqcDao.IQcIpqcDao
}

func NewQcIpqcServiceImpl(dao qcIpqcDao.IQcIpqcDao) qcIpqcService.IQcIpqcService {
	return &QcIpqcServiceImpl{dao: dao}
}

func (s *QcIpqcServiceImpl) GetQcIpqcList(ctx context.Context, req *qcIpqcApi.QcIpqcQueryReq) ([]*qcIpqcApi.QcIpqc, int64, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	models, total, err := s.dao.SelectQcIpqcList(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	apis := make([]*qcIpqcApi.QcIpqc, len(models))
	for i, m := range models {
		apis[i] = s.modelToApi(m)
	}
	return apis, total, nil
}

func (s *QcIpqcServiceImpl) GetQcIpqcById(ctx context.Context, id int64) (*qcIpqcApi.QcIpqc, error) {
	model, err := s.dao.SelectQcIpqcById(ctx, id)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}
	return s.modelToApi(model), nil
}

func (s *QcIpqcServiceImpl) CreateQcIpqc(ctx context.Context, req *qcIpqcApi.QcIpqcCreateReq) error {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return fmt.Errorf("上下文类型错误")
	}
	userName := "system"
	if ginCtx != nil {
		userName = ginCtx.GetString("userName")
	}
	model := &qcIpqcModels.QcIpqc{
		IpqcCode:              req.IpqcCode,
		IpqcName:              req.IpqcName,
		IpqcType:              req.IpqcType,
		TemplateId:            req.TemplateId,
		SourceDocId:           req.SourceDocId,
		SourceDocType:         req.SourceDocType,
		SourceDocCode:         req.SourceDocCode,
		SourceLineId:          req.SourceLineId,
		WorkorderId:           req.WorkorderId,
		WorkorderCode:         req.WorkorderCode,
		WorkorderName:         req.WorkorderName,
		TaskId:                req.TaskId,
		TaskCode:              req.TaskCode,
		TaskName:              req.TaskName,
		WorkstationId:         req.WorkstationId,
		WorkstationCode:       req.WorkstationCode,
		WorkstationName:       req.WorkstationName,
		ProcessId:             req.ProcessId,
		ProcessCode:           req.ProcessCode,
		ProcessName:           req.ProcessName,
		ItemId:                req.ItemId,
		ItemCode:              req.ItemCode,
		ItemName:              req.ItemName,
		Specification:         req.Specification,
		UnitOfMeasure:         req.UnitOfMeasure,
		UnitName:              req.UnitName,
		QuantityCheck:         req.QuantityCheck,
		QuantityUnqualified:   req.QuantityUnqualified,
		QuantityQualified:     req.QuantityQualified,
		QuantityLaborScrap:    req.QuantityLaborScrap,
		QuantityMaterialScrap: req.QuantityMaterialScrap,
		QuantityOtherScrap:    req.QuantityOtherScrap,
		CrRate:                req.CrRate,
		MajRate:               req.MajRate,
		MinRate:               req.MinRate,
		CrQuantity:            req.CrQuantity,
		MajQuantity:           req.MajQuantity,
		MinQuantity:           req.MinQuantity,
		CheckResult:           req.CheckResult,
		InspectDate:           req.InspectDate,
		Inspector:             req.Inspector,
		Status:                req.Status,
		Remark:                req.Remark,
		Attr1:                 req.Attr1,
		Attr2:                 req.Attr2,
		Attr3:                 req.Attr3,
		Attr4:                 req.Attr4,
		CreateBy:              &userName,
		UpdateBy:              &userName,
	}
	return s.dao.InsertQcIpqc(ctx, model)
}

func (s *QcIpqcServiceImpl) UpdateQcIpqc(ctx context.Context, req *qcIpqcApi.QcIpqcUpdateReq) error {
	model, err := s.dao.SelectQcIpqcById(ctx, req.IpqcId)
	if err != nil {
		return err
	}
	if model == nil {
		return fmt.Errorf("过程检验单不存在")
	}
	if req.IpqcName != nil {
		model.IpqcName = req.IpqcName
	}
	if req.IpqcType != nil {
		model.IpqcType = req.IpqcType
	}
	if req.TemplateId != nil {
		model.TemplateId = req.TemplateId
	}
	if req.Status != nil {
		model.Status = req.Status
	}
	if req.Remark != nil {
		model.Remark = req.Remark
	}
	if req.Attr1 != nil {
		model.Attr1 = req.Attr1
	}
	if req.Attr2 != nil {
		model.Attr2 = req.Attr2
	}
	if req.Attr3 != nil {
		model.Attr3 = req.Attr3
	}
	if req.Attr4 != nil {
		model.Attr4 = req.Attr4
	}
	return s.dao.UpdateQcIpqc(ctx, model)
}

func (s *QcIpqcServiceImpl) DeleteQcIpqcByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("ID列表不能为空")
	}
	return s.dao.DeleteQcIpqcByIds(ctx, ids)
}

func (s *QcIpqcServiceImpl) modelToApi(m *qcIpqcModels.QcIpqc) *qcIpqcApi.QcIpqc {
	return &qcIpqcApi.QcIpqc{
		IpqcId:                int64(m.ID),
		IpqcCode:              m.IpqcCode,
		IpqcName:              m.IpqcName,
		IpqcType:              m.IpqcType,
		TemplateId:            m.TemplateId,
		SourceDocId:           m.SourceDocId,
		SourceDocType:         m.SourceDocType,
		SourceDocCode:         m.SourceDocCode,
		SourceLineId:          m.SourceLineId,
		WorkorderId:           m.WorkorderId,
		WorkorderCode:         m.WorkorderCode,
		WorkorderName:         m.WorkorderName,
		TaskId:                m.TaskId,
		TaskCode:              m.TaskCode,
		TaskName:              m.TaskName,
		WorkstationId:         m.WorkstationId,
		WorkstationCode:       m.WorkstationCode,
		WorkstationName:       m.WorkstationName,
		ProcessId:             m.ProcessId,
		ProcessCode:           m.ProcessCode,
		ProcessName:           m.ProcessName,
		ItemId:                m.ItemId,
		ItemCode:              m.ItemCode,
		ItemName:              m.ItemName,
		Specification:         m.Specification,
		UnitOfMeasure:         m.UnitOfMeasure,
		UnitName:              m.UnitName,
		QuantityCheck:         m.QuantityCheck,
		QuantityUnqualified:   m.QuantityUnqualified,
		QuantityQualified:     m.QuantityQualified,
		QuantityLaborScrap:    m.QuantityLaborScrap,
		QuantityMaterialScrap: m.QuantityMaterialScrap,
		QuantityOtherScrap:    m.QuantityOtherScrap,
		CrRate:                m.CrRate,
		MajRate:               m.MajRate,
		MinRate:               m.MinRate,
		CrQuantity:            m.CrQuantity,
		MajQuantity:           m.MajQuantity,
		MinQuantity:           m.MinQuantity,
		CheckResult:           m.CheckResult,
		InspectDate:           m.InspectDate,
		Inspector:             m.Inspector,
		Status:                m.Status,
		Remark:                m.Remark,
		Attr1:                 m.Attr1,
		Attr2:                 m.Attr2,
		Attr3:                 m.Attr3,
		Attr4:                 m.Attr4,
		CreateBy:              derefStr(m.CreateBy),
		CreateTime:            m.CreatedAt,
		UpdateBy:              derefStr(m.UpdateBy),
		UpdateTime:            &m.UpdatedAt,
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
