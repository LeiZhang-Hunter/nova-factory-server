package qcTemplateService

import (
	"context"
	"nova-factory-server/app/business/qcTemplate/qcTemplateApi"
)

// IQcTemplateService 检测模板服务接口
type IQcTemplateService interface {
	// GetQcTemplateList 获取检测模板列表
	GetQcTemplateList(ctx context.Context, req *qcTemplateApi.QcTemplateQueryReq) ([]*qcTemplateApi.QcTemplate, int64, error)

	// GetQcTemplateById 根据ID获取检测模板
	GetQcTemplateById(ctx context.Context, templateId int64) (*qcTemplateApi.QcTemplate, error)

	// CreateQcTemplate 创建检测模板
	CreateQcTemplate(ctx context.Context, req *qcTemplateApi.QcTemplateCreateReq) error

	// UpdateQcTemplate 更新检测模板
	UpdateQcTemplate(ctx context.Context, req *qcTemplateApi.QcTemplateUpdateReq) error

	// DeleteQcTemplateByIds 批量删除检测模板
	DeleteQcTemplateByIds(ctx context.Context, templateIds []int64) error
}
