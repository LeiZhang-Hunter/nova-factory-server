package qcIpqcDaoImpl

import (
	"context"
	"nova-factory-server/app/business/qcIpqc/qcIpqcApi"
	"nova-factory-server/app/business/qcIpqc/qcIpqcDao"
	"nova-factory-server/app/business/qcIpqc/qcIpqcModels"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type QcIpqcDaoImpl struct {
	db *gorm.DB
}

func NewQcIpqcDaoImpl(db *gorm.DB) qcIpqcDao.IQcIpqcDao {
	return &QcIpqcDaoImpl{db: db}
}

func (dao *QcIpqcDaoImpl) SelectQcIpqcList(ctx context.Context, req *qcIpqcApi.QcIpqcQueryReq) ([]*qcIpqcModels.QcIpqc, int64, error) {
	var result []*qcIpqcModels.QcIpqc
	var total int64
	query := dao.db.WithContext(ctx).Model(&qcIpqcModels.QcIpqc{})
	if req.IpqcCode != nil && *req.IpqcCode != "" {
		query = query.Where("ipqc_code = ?", *req.IpqcCode)
	}
	if req.IpqcName != nil && *req.IpqcName != "" {
		query = query.Where("ipqc_name LIKE ?", "%"+*req.IpqcName+"%")
	}
	if req.IpqcType != nil && *req.IpqcType != "" {
		query = query.Where("ipqc_type = ?", *req.IpqcType)
	}
	if req.TemplateId != nil {
		query = query.Where("template_id = ?", *req.TemplateId)
	}
	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	err := query.Count(&total).Error
	if err != nil {
		zap.L().Error("查询过程检验单总数失败", zap.Error(err))
		return nil, 0, err
	}
	err = query.Order("created_at DESC").Limit(req.PageSize).Offset((req.PageNum - 1) * req.PageSize).Find(&result).Error
	if err != nil {
		zap.L().Error("查询过程检验单列表失败", zap.Error(err))
		return nil, 0, err
	}
	return result, total, nil
}

func (dao *QcIpqcDaoImpl) SelectQcIpqcById(ctx context.Context, id int64) (*qcIpqcModels.QcIpqc, error) {
	var ipqc qcIpqcModels.QcIpqc
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&ipqc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		zap.L().Error("根据ID查询过程检验单失败", zap.Error(err))
		return nil, err
	}
	return &ipqc, nil
}

func (dao *QcIpqcDaoImpl) InsertQcIpqc(ctx context.Context, ipqc *qcIpqcModels.QcIpqc) error {
	err := dao.db.WithContext(ctx).Create(ipqc).Error
	if err != nil {
		zap.L().Error("新增过程检验单失败", zap.Error(err))
		return err
	}
	return nil
}

func (dao *QcIpqcDaoImpl) UpdateQcIpqc(ctx context.Context, ipqc *qcIpqcModels.QcIpqc) error {
	err := dao.db.WithContext(ctx).Save(ipqc).Error
	if err != nil {
		zap.L().Error("更新过程检验单失败", zap.Error(err))
		return err
	}
	return nil
}

func (dao *QcIpqcDaoImpl) DeleteQcIpqcByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	err := dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&qcIpqcModels.QcIpqc{}).Error
	if err != nil {
		zap.L().Error("批量删除过程检验单失败", zap.Error(err))
		return err
	}
	return nil
}
