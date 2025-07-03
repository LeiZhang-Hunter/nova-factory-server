package qcRqcDaoImpl

import (
	"context"

	"nova-factory-server/app/business/qcRqc/qcRqcApi"
	"nova-factory-server/app/business/qcRqc/qcRqcDao"
	"nova-factory-server/app/business/qcRqc/qcRqcModels"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// QcRqcDaoImpl 退料检验单数据访问实现
type QcRqcDaoImpl struct {
	db *gorm.DB
}

// NewQcRqcDaoImpl 创建退料检验单DAO实现
func NewQcRqcDaoImpl(db *gorm.DB) qcRqcDao.IQcRqcDao {
	return &QcRqcDaoImpl{db: db}
}

// SelectQcRqcList 查询退料检验单列表
func (dao *QcRqcDaoImpl) SelectQcRqcList(ctx context.Context, req *qcRqcApi.QcRqcQueryReq) ([]*qcRqcModels.QcRqc, int64, error) {
	var result []*qcRqcModels.QcRqc
	var total int64

	// 构建查询
	query := dao.db.WithContext(ctx).Model(&qcRqcModels.QcRqc{})

	// 应用查询条件
	query = dao.applyWhereConditions(query, req)

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("查询退料检验单总数失败", zap.Error(err))
		return nil, 0, err
	}

	// 查询数据
	err = query.Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int((req.PageNum - 1) * req.PageSize)).
		Find(&result).Error

	if err != nil {
		zap.L().Error("查询退料检验单列表失败", zap.Error(err))
		return nil, 0, err
	}

	return result, total, nil
}

// SelectQcRqcById 根据ID查询退料检验单
func (dao *QcRqcDaoImpl) SelectQcRqcById(ctx context.Context, id int64) (*qcRqcModels.QcRqc, error) {
	var rqc qcRqcModels.QcRqc

	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&rqc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		zap.L().Error("根据ID查询退料检验单失败", zap.Error(err))
		return nil, err
	}

	return &rqc, nil
}

// InsertQcRqc 新增退料检验单
func (dao *QcRqcDaoImpl) InsertQcRqc(ctx context.Context, rqc *qcRqcModels.QcRqc) error {
	err := dao.db.WithContext(ctx).Create(rqc).Error
	if err != nil {
		zap.L().Error("新增退料检验单失败", zap.Error(err))
		return err
	}

	return nil
}

// UpdateQcRqc 更新退料检验单
func (dao *QcRqcDaoImpl) UpdateQcRqc(ctx context.Context, rqc *qcRqcModels.QcRqc) error {
	err := dao.db.WithContext(ctx).Save(rqc).Error
	if err != nil {
		zap.L().Error("更新退料检验单失败", zap.Error(err))
		return err
	}

	return nil
}

// DeleteQcRqcByIds 批量删除退料检验单
func (dao *QcRqcDaoImpl) DeleteQcRqcByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	err := dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&qcRqcModels.QcRqc{}).Error
	if err != nil {
		zap.L().Error("批量删除退料检验单失败", zap.Error(err))
		return err
	}

	return nil
}

// applyWhereConditions 应用查询条件
func (dao *QcRqcDaoImpl) applyWhereConditions(query *gorm.DB, req *qcRqcApi.QcRqcQueryReq) *gorm.DB {
	if req.RqcCode != nil && *req.RqcCode != "" {
		query = query.Where("rqc_code = ?", *req.RqcCode)
	}
	if req.RqcName != nil && *req.RqcName != "" {
		query = query.Where("rqc_name LIKE ?", "%"+*req.RqcName+"%")
	}
	if req.TemplateId != nil {
		query = query.Where("template_id = ?", *req.TemplateId)
	}
	if req.SourceDocId != nil {
		query = query.Where("source_doc_id = ?", *req.SourceDocId)
	}
	if req.SourceDocType != nil && *req.SourceDocType != "" {
		query = query.Where("source_doc_type = ?", *req.SourceDocType)
	}
	if req.SourceDocCode != nil && *req.SourceDocCode != "" {
		query = query.Where("source_doc_code = ?", *req.SourceDocCode)
	}
	if req.SourceLineId != nil {
		query = query.Where("source_line_id = ?", *req.SourceLineId)
	}
	if req.ItemId != nil {
		query = query.Where("item_id = ?", *req.ItemId)
	}
	if req.ItemCode != nil && *req.ItemCode != "" {
		query = query.Where("item_code = ?", *req.ItemCode)
	}
	if req.ItemName != nil && *req.ItemName != "" {
		query = query.Where("item_name LIKE ?", "%"+*req.ItemName+"%")
	}
	if req.Specification != nil && *req.Specification != "" {
		query = query.Where("specification LIKE ?", "%"+*req.Specification+"%")
	}
	if req.UnitOfMeasure != nil && *req.UnitOfMeasure != "" {
		query = query.Where("unit_of_measure = ?", *req.UnitOfMeasure)
	}
	if req.UnitName != nil && *req.UnitName != "" {
		query = query.Where("unit_name LIKE ?", "%"+*req.UnitName+"%")
	}
	if req.BatchId != nil {
		query = query.Where("batch_id = ?", *req.BatchId)
	}
	if req.BatchCode != nil && *req.BatchCode != "" {
		query = query.Where("batch_code = ?", *req.BatchCode)
	}
	if req.QuantityCheck != nil {
		query = query.Where("quantity_check = ?", *req.QuantityCheck)
	}
	if req.QuantityUnqualified != nil {
		query = query.Where("quantity_unqualified = ?", *req.QuantityUnqualified)
	}
	if req.QuantityQualified != nil {
		query = query.Where("quantity_qualified = ?", *req.QuantityQualified)
	}
	if req.CheckResult != nil && *req.CheckResult != "" {
		query = query.Where("check_result = ?", *req.CheckResult)
	}
	if req.InspectDate != nil {
		query = query.Where("inspect_date = ?", *req.InspectDate)
	}
	if req.UserId != nil {
		query = query.Where("user_id = ?", *req.UserId)
	}
	if req.UserName != nil && *req.UserName != "" {
		query = query.Where("user_name LIKE ?", "%"+*req.UserName+"%")
	}
	if req.NickName != nil && *req.NickName != "" {
		query = query.Where("nick_name LIKE ?", "%"+*req.NickName+"%")
	}
	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	if req.CreateBy != nil && *req.CreateBy != "" {
		query = query.Where("create_by = ?", *req.CreateBy)
	}

	return query
}
