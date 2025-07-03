package qcOqcDaoImpl

import (
	"context"

	"nova-factory-server/app/business/qcOqc/qcOqcApi"
	"nova-factory-server/app/business/qcOqc/qcOqcDao"
	"nova-factory-server/app/business/qcOqc/qcOqcModels"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// QcOqcDaoImpl 出货检验单数据访问实现
type QcOqcDaoImpl struct {
	db *gorm.DB
}

// NewQcOqcDaoImpl 创建出货检验单DAO实现
func NewQcOqcDaoImpl(db *gorm.DB) qcOqcDao.IQcOqcDao {
	return &QcOqcDaoImpl{db: db}
}

// SelectQcOqcList 查询出货检验单列表
func (dao *QcOqcDaoImpl) SelectQcOqcList(ctx context.Context, req *qcOqcApi.QcOqcQueryReq) ([]*qcOqcModels.QcOqc, int64, error) {
	var result []*qcOqcModels.QcOqc
	var total int64

	// 构建查询
	query := dao.db.WithContext(ctx).Model(&qcOqcModels.QcOqc{})

	// 应用查询条件
	query = dao.applyWhereConditions(query, req)

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("查询出货检验单总数失败", zap.Error(err))
		return nil, 0, err
	}

	// 查询数据
	err = query.Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int((req.PageNum - 1) * req.PageSize)).
		Find(&result).Error

	if err != nil {
		zap.L().Error("查询出货检验单列表失败", zap.Error(err))
		return nil, 0, err
	}

	return result, total, nil
}

// SelectQcOqcById 根据ID查询出货检验单
func (dao *QcOqcDaoImpl) SelectQcOqcById(ctx context.Context, id int64) (*qcOqcModels.QcOqc, error) {
	var oqc qcOqcModels.QcOqc

	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&oqc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		zap.L().Error("根据ID查询出货检验单失败", zap.Error(err))
		return nil, err
	}

	return &oqc, nil
}

// InsertQcOqc 新增出货检验单
func (dao *QcOqcDaoImpl) InsertQcOqc(ctx context.Context, oqc *qcOqcModels.QcOqc) error {
	err := dao.db.WithContext(ctx).Create(oqc).Error
	if err != nil {
		zap.L().Error("新增出货检验单失败", zap.Error(err))
		return err
	}

	return nil
}

// UpdateQcOqc 更新出货检验单
func (dao *QcOqcDaoImpl) UpdateQcOqc(ctx context.Context, oqc *qcOqcModels.QcOqc) error {
	err := dao.db.WithContext(ctx).Save(oqc).Error
	if err != nil {
		zap.L().Error("更新出货检验单失败", zap.Error(err))
		return err
	}

	return nil
}

// DeleteQcOqcByIds 批量删除出货检验单
func (dao *QcOqcDaoImpl) DeleteQcOqcByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	err := dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&qcOqcModels.QcOqc{}).Error
	if err != nil {
		zap.L().Error("批量删除出货检验单失败", zap.Error(err))
		return err
	}

	return nil
}

// applyWhereConditions 应用查询条件
func (dao *QcOqcDaoImpl) applyWhereConditions(query *gorm.DB, req *qcOqcApi.QcOqcQueryReq) *gorm.DB {
	if req.OqcCode != nil && *req.OqcCode != "" {
		query = query.Where("oqc_code = ?", *req.OqcCode)
	}
	if req.OqcName != nil && *req.OqcName != "" {
		query = query.Where("oqc_name LIKE ?", "%"+*req.OqcName+"%")
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
	if req.ClientId != nil {
		query = query.Where("client_id = ?", *req.ClientId)
	}
	if req.ClientCode != nil && *req.ClientCode != "" {
		query = query.Where("client_code = ?", *req.ClientCode)
	}
	if req.ClientName != nil && *req.ClientName != "" {
		query = query.Where("client_name LIKE ?", "%"+*req.ClientName+"%")
	}
	if req.BatchCode != nil && *req.BatchCode != "" {
		query = query.Where("batch_code = ?", *req.BatchCode)
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
	if req.QuantityMinCheck != nil {
		query = query.Where("quantity_min_check = ?", *req.QuantityMinCheck)
	}
	if req.QuantityMaxUnqualified != nil {
		query = query.Where("quantity_max_unqualified = ?", *req.QuantityMaxUnqualified)
	}
	if req.QuantityOut != nil {
		query = query.Where("quantity_out = ?", *req.QuantityOut)
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
	if req.CrRate != nil {
		query = query.Where("cr_rate = ?", *req.CrRate)
	}
	if req.MajRate != nil {
		query = query.Where("maj_rate = ?", *req.MajRate)
	}
	if req.MinRate != nil {
		query = query.Where("min_rate = ?", *req.MinRate)
	}
	if req.CrQuantity != nil {
		query = query.Where("cr_quantity = ?", *req.CrQuantity)
	}
	if req.MajQuantity != nil {
		query = query.Where("maj_quantity = ?", *req.MajQuantity)
	}
	if req.MinQuantity != nil {
		query = query.Where("min_quantity = ?", *req.MinQuantity)
	}
	if req.CheckResult != nil && *req.CheckResult != "" {
		query = query.Where("check_result = ?", *req.CheckResult)
	}
	if req.OutDate != nil {
		query = query.Where("out_date = ?", *req.OutDate)
	}
	if req.InspectDate != nil {
		query = query.Where("inspect_date = ?", *req.InspectDate)
	}
	if req.Inspector != nil && *req.Inspector != "" {
		query = query.Where("inspector = ?", *req.Inspector)
	}
	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	if req.CreateBy != nil && *req.CreateBy != "" {
		query = query.Where("create_by = ?", *req.CreateBy)
	}

	return query
}
