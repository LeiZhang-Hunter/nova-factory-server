package productDaoImpl

import (
	"context"
	"errors"
	"nova-factory-server/app/business/admin/product/productdao"
	"nova-factory-server/app/business/admin/product/productmodels"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type sysProductLaboratoryDao struct {
	db    *gorm.DB
	table string
}

func NewSysProductLaboratoryDao(db *gorm.DB) productdao.ISysProductLaboratoryDao {
	return &sysProductLaboratoryDao{
		db:    db,
		table: "sys_product_laboratory",
	}
}

func (dao *sysProductLaboratoryDao) SelectLaboratoryList(ctx *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (list *productmodels.SysProductLaboratoryList, err error) {
	if dql == nil {
		dql = &productmodels.SysProductLaboratoryDQL{}
	}
	query := dao.db.Table(dao.table)

	query = query.Where("type = ?", dql.Type)

	if dql.Material != "" {
		query = query.Where("material = ?", dql.Material)
	}
	if dql.Contact != "" {
		query = query.Where("contact LIKE ?", "%"+dql.Contact+"%")
	}
	query = baizeContext.GetGormDataScope(ctx, query)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &productmodels.SysProductLaboratoryList{
			Rows:  make([]*productmodels.SysProductLaboratory, 0),
			Total: 0,
		}, nil
	}

	size := 0
	if dql.Size <= 0 {
		size = 20
	} else {
		size = int(dql.Size)
	}
	offset := 0
	if dql.Page <= 0 {
		dql.Page = 1
	} else {
		offset = int((dql.Page - 1) * dql.Size)
	}

	dto := make([]*productmodels.SysProductLaboratory, 0)
	ret := query.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &productmodels.SysProductLaboratoryList{
			Rows:  make([]*productmodels.SysProductLaboratory, 0),
			Total: 0,
		}, ret.Error
	}
	return &productmodels.SysProductLaboratoryList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (dao *sysProductLaboratoryDao) SelectLaboratoryById(ctx context.Context, id int64) (*productmodels.SysProductLaboratoryVo, error) {
	data := new(productmodels.SysProductLaboratoryVo)
	err := dao.db.
		Table(dao.table).
		Where("id = ?", id).
		First(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return data, err
	}
	return data, nil
}

func (dao *sysProductLaboratoryDao) Set(c *gin.Context, data *productmodels.SysProductLaboratoryVo) (*productmodels.SysProductLaboratory, error) {
	value := productmodels.ToSysProductLaboratory(data)
	if value.Id == 0 {
		value.SetCreateBy(baizeContext.GetUserId(c))
		value.Contact = baizeContext.GetUserName(c)
		value.Id = snowflake.GenID()
		value.DeptId = baizeContext.GetDeptId(c)
		ret := dao.db.Table(dao.table).Create(&value)
		return value, ret.Error
	} else {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		value.Contact = baizeContext.GetUserName(c)
		ret := dao.db.Table(dao.table).Where("id = ?", value.Id).Updates(&value)
		return value, ret.Error
	}
}

func (dao *sysProductLaboratoryDao) DeleteLaboratoryByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	if err := dao.db.WithContext(ctx).
		Table(dao.table).
		Where("id IN ?", ids).
		Delete(&productmodels.SysProductLaboratoryVo{}).Error; err != nil {
		return err
	}
	return nil
}

func (dao *sysProductLaboratoryDao) SelectUserLaboratoryList(ctx *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (list *productmodels.SysProductLaboratoryList, err error) {
	if dql == nil {
		dql = &productmodels.SysProductLaboratoryDQL{}
	}
	query := dao.db.Table(dao.table)
	userId := baizeContext.GetUserId(ctx)
	query = query.Where("create_by = ?", userId)
	query = baizeContext.GetGormDataScope(ctx, query)

	if dql.BeginTime != "" {
		query = query.Where("create_time > ?", dql.BeginTime)
	}

	if dql.EndTime != "" {
		query = query.Where("create_time <= ?", dql.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &productmodels.SysProductLaboratoryList{
			Rows:  make([]*productmodels.SysProductLaboratory, 0),
			Total: 0,
		}, nil
	}

	size := 0
	if dql.Size <= 0 {
		size = 20
	} else {
		size = int(dql.Size)
	}
	offset := 0
	if dql.Page <= 0 {
		dql.Page = 1
	} else {
		offset = int((dql.Page - 1) * dql.Size)
	}

	dto := make([]*productmodels.SysProductLaboratory, 0)
	ret := query.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &productmodels.SysProductLaboratoryList{
			Rows:  make([]*productmodels.SysProductLaboratory, 0),
			Total: 0,
		}, ret.Error
	}
	return &productmodels.SysProductLaboratoryList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (dao *sysProductLaboratoryDao) FirstLaboratoryInfo(ctx *gin.Context, req *productmodels.SysProductLaboratoryInfoDQL) (*productmodels.SysProductLaboratory, error) {
	query := dao.db.Table(dao.table)
	var dto productmodels.SysProductLaboratory
	ret := query.Where("type = ?", req.Type).Order("create_time desc").Limit(1).Find(&dto)
	if ret.Error != nil {
		return &dto, ret.Error
	}
	return &dto, nil
}

func (dao *sysProductLaboratoryDao) FirstLaboratoryList(ctx *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (*productmodels.SysProductLaboratoryList, error) {
	if dql == nil {
		dql = &productmodels.SysProductLaboratoryDQL{}
	}
	query := dao.db.Table(dao.table)
	query.Where("type = ?", dql.Type)

	if dql.Material != "" {
		query.Where("material = ?", dql.Material)
	}

	size := 0
	if dql.Size <= 0 {
		size = 20
	} else {
		size = int(dql.Size)
	}
	offset := 0
	if dql.Page <= 0 {
		dql.Page = 1
	} else {
		offset = int((dql.Page - 1) * dql.Size)
	}

	query = query.Where("type = ?", dql.Type)

	dto := make([]*productmodels.SysProductLaboratory, 0)
	ret := query.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &productmodels.SysProductLaboratoryList{
			Rows:  make([]*productmodels.SysProductLaboratory, 0),
			Total: 0,
		}, ret.Error
	}
	return &productmodels.SysProductLaboratoryList{
		Rows:  dto,
		Total: int64(len(dto)),
	}, nil
}
