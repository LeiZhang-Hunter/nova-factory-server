package qcTemplateDao

import (
	"context"
	"nova-factory-server/app/business/qcTemplate/qcTemplateApi"
	"nova-factory-server/app/business/qcTemplate/qcTemplateModels"
)

// IQcTemplateDao 检测模板数据访问接口
type IQcTemplateDao interface {
	// SelectQcTemplateList 查询检测模板列表
	SelectQcTemplateList(ctx context.Context, req *qcTemplateApi.QcTemplateQueryReq) ([]*qcTemplateModels.QcTemplate, int64, error)

	// SelectQcTemplateById 根据ID查询检测模板
	SelectQcTemplateById(ctx context.Context, id int64) (*qcTemplateModels.QcTemplate, error)

	// InsertQcTemplate 新增检测模板
	InsertQcTemplate(ctx context.Context, template *qcTemplateModels.QcTemplate) error

	// UpdateQcTemplate 更新检测模板
	UpdateQcTemplate(ctx context.Context, template *qcTemplateModels.QcTemplate) error

	// DeleteQcTemplateByIds 批量删除检测模板
	DeleteQcTemplateByIds(ctx context.Context, ids []int64) error
}
