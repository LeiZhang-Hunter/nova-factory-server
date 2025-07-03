package qcTemplateDaoImpl

import (
	"context"

	"nova-factory-server/app/business/qcTemplate/qcTemplateApi"
	"nova-factory-server/app/business/qcTemplate/qcTemplateDao"
	"nova-factory-server/app/business/qcTemplate/qcTemplateModels"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// QcTemplateDaoImpl 检测模板数据访问实现
type QcTemplateDaoImpl struct {
	db *gorm.DB
}

// NewQcTemplateDaoImpl 创建检测模板DAO实现
func NewQcTemplateDaoImpl(db *gorm.DB) qcTemplateDao.IQcTemplateDao {
	return &QcTemplateDaoImpl{db: db}
}

// SelectQcTemplateList 查询检测模板列表
func (dao *QcTemplateDaoImpl) SelectQcTemplateList(ctx context.Context, req *qcTemplateApi.QcTemplateQueryReq) ([]*qcTemplateModels.QcTemplate, int64, error) {
	var result []*qcTemplateModels.QcTemplate
	var total int64

	// 构建查询
	query := dao.db.WithContext(ctx).Model(&qcTemplateModels.QcTemplate{})

	// 应用查询条件
	query = dao.applyWhereConditions(query, req)

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("查询检测模板总数失败", zap.Error(err))
		return nil, 0, err
	}

	// 查询数据
	err = query.Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int((req.PageNum - 1) * req.PageSize)).
		Find(&result).Error

	if err != nil {
		zap.L().Error("查询检测模板列表失败", zap.Error(err))
		return nil, 0, err
	}

	return result, total, nil
}

// SelectQcTemplateById 根据ID查询检测模板
func (dao *QcTemplateDaoImpl) SelectQcTemplateById(ctx context.Context, id int64) (*qcTemplateModels.QcTemplate, error) {
	var template qcTemplateModels.QcTemplate

	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&template).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		zap.L().Error("根据ID查询检测模板失败", zap.Error(err))
		return nil, err
	}

	return &template, nil
}

// InsertQcTemplate 新增检测模板
func (dao *QcTemplateDaoImpl) InsertQcTemplate(ctx context.Context, template *qcTemplateModels.QcTemplate) error {
	err := dao.db.WithContext(ctx).Create(template).Error
	if err != nil {
		zap.L().Error("新增检测模板失败", zap.Error(err))
		return err
	}

	return nil
}

// UpdateQcTemplate 更新检测模板
func (dao *QcTemplateDaoImpl) UpdateQcTemplate(ctx context.Context, template *qcTemplateModels.QcTemplate) error {
	err := dao.db.WithContext(ctx).Save(template).Error
	if err != nil {
		zap.L().Error("更新检测模板失败", zap.Error(err))
		return err
	}

	return nil
}

// DeleteQcTemplateByIds 批量删除检测模板
func (dao *QcTemplateDaoImpl) DeleteQcTemplateByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	err := dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&qcTemplateModels.QcTemplate{}).Error
	if err != nil {
		zap.L().Error("批量删除检测模板失败", zap.Error(err))
		return err
	}

	return nil
}

// applyWhereConditions 应用查询条件
func (dao *QcTemplateDaoImpl) applyWhereConditions(query *gorm.DB, req *qcTemplateApi.QcTemplateQueryReq) *gorm.DB {
	if req.TemplateCode != nil && *req.TemplateCode != "" {
		query = query.Where("template_code LIKE ?", "%"+*req.TemplateCode+"%")
	}
	if req.TemplateName != nil && *req.TemplateName != "" {
		query = query.Where("template_name LIKE ?", "%"+*req.TemplateName+"%")
	}
	if req.QcTypes != nil && *req.QcTypes != "" {
		query = query.Where("qc_types LIKE ?", "%"+*req.QcTypes+"%")
	}
	if req.EnableFlag != nil && *req.EnableFlag != "" {
		query = query.Where("enable_flag = ?", *req.EnableFlag)
	}
	if req.CreateBy != nil && *req.CreateBy != "" {
		query = query.Where("create_by = ?", *req.CreateBy)
	}

	return query
}
