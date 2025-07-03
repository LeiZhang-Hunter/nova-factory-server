package qcTemplateServiceImpl

import (
	"context"
	"fmt"
	"time"

	"nova-factory-server/app/business/qcTemplate/qcTemplateApi"
	"nova-factory-server/app/business/qcTemplate/qcTemplateDao"
	"nova-factory-server/app/business/qcTemplate/qcTemplateModels"
	"nova-factory-server/app/business/qcTemplate/qcTemplateService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QcTemplateServiceImpl 检测模板服务实现
type QcTemplateServiceImpl struct {
	qcTemplateDao qcTemplateDao.IQcTemplateDao
}

// NewQcTemplateServiceImpl 创建检测模板服务实现
func NewQcTemplateServiceImpl(qcTemplateDao qcTemplateDao.IQcTemplateDao) qcTemplateService.IQcTemplateService {
	return &QcTemplateServiceImpl{qcTemplateDao: qcTemplateDao}
}

// GetQcTemplateList 获取检测模板列表
func (s *QcTemplateServiceImpl) GetQcTemplateList(ctx context.Context, req *qcTemplateApi.QcTemplateQueryReq) ([]*qcTemplateApi.QcTemplate, int64, error) {
	// 设置默认分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用DAO查询数据
	models, total, err := s.qcTemplateDao.SelectQcTemplateList(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API结构体
	apis := make([]*qcTemplateApi.QcTemplate, len(models))
	for i, model := range models {
		apis[i] = s.convertModelToApi(model)
	}

	return apis, total, nil
}

// GetQcTemplateById 根据ID获取检测模板
func (s *QcTemplateServiceImpl) GetQcTemplateById(ctx context.Context, templateId int64) (*qcTemplateApi.QcTemplate, error) {
	if templateId <= 0 {
		return nil, fmt.Errorf("检测模板ID不能为空")
	}

	model, err := s.qcTemplateDao.SelectQcTemplateById(ctx, templateId)
	if err != nil {
		return nil, err
	}

	if model == nil {
		return nil, nil
	}

	return s.convertModelToApi(model), nil
}

// CreateQcTemplate 创建检测模板
func (s *QcTemplateServiceImpl) CreateQcTemplate(ctx context.Context, req *qcTemplateApi.QcTemplateCreateReq) error {
	// 获取当前用户
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return fmt.Errorf("上下文类型错误")
	}
	userName := baizeContext.GetUserName(ginCtx)
	if userName == "" {
		return fmt.Errorf("获取用户信息失败")
	}

	// 生成模板编号
	templateCode := fmt.Sprintf("TPL%s%04d", time.Now().Format("20060102"), snowflake.GenID()%10000)

	// 构建检测模板model对象
	template := &qcTemplateModels.QcTemplate{
		TemplateCode: templateCode,
		TemplateName: req.TemplateName,
		QcTypes:      req.QcTypes,
		EnableFlag:   req.EnableFlag,
		Remark:       req.Remark,
		Attr1:        req.Attr1,
		Attr2:        req.Attr2,
		Attr3:        req.Attr3,
		Attr4:        req.Attr4,
		CreateBy:     &userName,
		UpdateBy:     &userName,
		CreateById:   baizeContext.GetUserId(ginCtx),
		UpdateById:   baizeContext.GetUserId(ginCtx),
	}

	// 设置默认启用状态
	if template.EnableFlag == nil {
		defaultFlag := "Y"
		template.EnableFlag = &defaultFlag
	}

	return s.qcTemplateDao.InsertQcTemplate(ctx, template)
}

// UpdateQcTemplate 更新检测模板
func (s *QcTemplateServiceImpl) UpdateQcTemplate(ctx context.Context, req *qcTemplateApi.QcTemplateUpdateReq) error {
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
	existingTemplate, err := s.qcTemplateDao.SelectQcTemplateById(ctx, req.TemplateId)
	if err != nil {
		zap.L().Error("查询检测模板失败", zap.Error(err))
		return err
	}
	if existingTemplate == nil {
		return fmt.Errorf("检测模板不存在")
	}

	// 更新字段
	if req.TemplateCode != nil {
		existingTemplate.TemplateCode = *req.TemplateCode
	}
	if req.TemplateName != nil {
		existingTemplate.TemplateName = *req.TemplateName
	}
	if req.QcTypes != nil {
		existingTemplate.QcTypes = *req.QcTypes
	}
	if req.EnableFlag != nil {
		existingTemplate.EnableFlag = req.EnableFlag
	}
	if req.Remark != nil {
		existingTemplate.Remark = req.Remark
	}
	if req.Attr1 != nil {
		existingTemplate.Attr1 = req.Attr1
	}
	if req.Attr2 != nil {
		existingTemplate.Attr2 = req.Attr2
	}
	if req.Attr3 != nil {
		existingTemplate.Attr3 = req.Attr3
	}
	if req.Attr4 != nil {
		existingTemplate.Attr4 = req.Attr4
	}

	// 更新用户信息
	existingTemplate.UpdateBy = &userName
	existingTemplate.UpdateById = baizeContext.GetUserId(ginCtx)

	return s.qcTemplateDao.UpdateQcTemplate(ctx, existingTemplate)
}

// DeleteQcTemplateByIds 批量删除检测模板
func (s *QcTemplateServiceImpl) DeleteQcTemplateByIds(ctx context.Context, templateIds []int64) error {
	if len(templateIds) == 0 {
		return fmt.Errorf("检测模板ID列表不能为空")
	}

	return s.qcTemplateDao.DeleteQcTemplateByIds(ctx, templateIds)
}

// convertModelToApi 将model结构体转换为API结构体
func (s *QcTemplateServiceImpl) convertModelToApi(model *qcTemplateModels.QcTemplate) *qcTemplateApi.QcTemplate {
	return &qcTemplateApi.QcTemplate{
		TemplateId:   int64(model.ID),
		TemplateCode: model.TemplateCode,
		TemplateName: model.TemplateName,
		QcTypes:      model.QcTypes,
		EnableFlag:   model.EnableFlag,
		Remark:       model.Remark,
		Attr1:        model.Attr1,
		Attr2:        model.Attr2,
		Attr3:        model.Attr3,
		Attr4:        model.Attr4,
		CreateBy:     s.stringPtrToString(model.CreateBy),
		CreateTime:   model.CreatedAt,
		UpdateBy:     s.stringPtrToString(model.UpdateBy),
		UpdateTime:   &model.UpdatedAt,
	}
}

// stringPtrToString 将字符串指针转换为字符串
func (s *QcTemplateServiceImpl) stringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
