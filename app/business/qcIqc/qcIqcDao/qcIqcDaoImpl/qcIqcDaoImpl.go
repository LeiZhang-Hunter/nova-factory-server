package qcIqcDaoImpl

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIqc/qcIqcApi"
	"nova-factory-server/app/business/qcIqc/qcIqcDao"
	"nova-factory-server/app/business/qcIqc/qcIqcModel"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QcIqcDaoImpl 来料检验单数据访问实现
type QcIqcDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewQcIqcDaoImpl 创建来料检验单数据访问实现
func NewQcIqcDaoImpl(db *gorm.DB) qcIqcDao.IQcIqcDao {
	return &QcIqcDaoImpl{
		db:        db,
		tableName: "qc_iqc",
	}
}

// List 查询来料检验单列表
func (dao *QcIqcDaoImpl) List(c *gin.Context, req *qcIqcApi.QcIqcQueryReq) (int64, []*qcIqcModel.QcIqc, error) {
	db := dao.db.Model(&qcIqcModel.QcIqc{})

	// 添加查询条件
	if req.IqcCode != nil && *req.IqcCode != "" {
		db = db.Where("iqc_code = ?", *req.IqcCode)
	}
	if req.IqcName != nil && *req.IqcName != "" {
		db = db.Where("iqc_name LIKE ?", "%"+*req.IqcName+"%")
	}
	if req.TemplateId != nil {
		db = db.Where("template_id = ?", *req.TemplateId)
	}
	if req.VendorId != nil {
		db = db.Where("vendor_id = ?", *req.VendorId)
	}
	if req.VendorCode != nil && *req.VendorCode != "" {
		db = db.Where("vendor_code = ?", *req.VendorCode)
	}
	if req.VendorName != nil && *req.VendorName != "" {
		db = db.Where("vendor_name LIKE ?", "%"+*req.VendorName+"%")
	}
	if req.ItemId != nil {
		db = db.Where("item_id = ?", *req.ItemId)
	}
	if req.ItemCode != nil && *req.ItemCode != "" {
		db = db.Where("item_code = ?", *req.ItemCode)
	}
	if req.ItemName != nil && *req.ItemName != "" {
		db = db.Where("item_name LIKE ?", "%"+*req.ItemName+"%")
	}
	if req.Status != nil && *req.Status != "" {
		db = db.Where("status = ?", *req.Status)
	}

	size := 0
	if req == nil || req.PageSize <= 0 {
		size = 20
	} else {
		size = int(req.PageSize)
	}
	offset := 0
	if req == nil || req.PageNum <= 0 {
		req.PageNum = 1
	} else {
		offset = int((req.PageNum - 1) * req.PageSize)
	}

	var dto []*qcIqcModel.QcIqc
	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	ret = db.Offset(offset).Order("id desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	return total, dto, nil
}

// Create 新增来料检验单
func (dao *QcIqcDaoImpl) Create(c *gin.Context, req *qcIqcApi.QcIqcCreateReq) (*baize.EmptyResponse, error) {
	iqc := &qcIqcModel.QcIqc{
		IqcName:                req.IqcName,
		TemplateId:             req.TemplateId,
		SourceDocId:            req.SourceDocId,
		SourceDocType:          req.SourceDocType,
		SourceDocCode:          req.SourceDocCode,
		SourceLineId:           req.SourceLineId,
		VendorId:               req.VendorId,
		VendorCode:             req.VendorCode,
		VendorName:             req.VendorName,
		VendorNick:             req.VendorNick,
		VendorBatch:            req.VendorBatch,
		ItemId:                 req.ItemId,
		ItemCode:               req.ItemCode,
		ItemName:               req.ItemName,
		Specification:          req.Specification,
		UnitOfMeasure:          req.UnitOfMeasure,
		UnitName:               req.UnitName,
		QuantityMinCheck:       req.QuantityMinCheck,
		QuantityMaxUnqualified: req.QuantityMaxUnqualified,
		QuantityRecived:        req.QuantityRecived,
		QuantityCheck:          req.QuantityCheck,
		QuantityQualified:      req.QuantityQualified,
		QuantityUnqualified:    req.QuantityUnqualified,
		CrRate:                 req.CrRate,
		MajRate:                req.MajRate,
		MinRate:                req.MinRate,
		CrQuantity:             req.CrQuantity,
		MajQuantity:            req.MajQuantity,
		MinQuantity:            req.MinQuantity,
		CheckResult:            req.CheckResult,
		ReciveDate:             req.ReciveDate,
		InspectDate:            req.InspectDate,
		Inspector:              req.Inspector,
		Status:                 req.Status,
		Remark:                 req.Remark,
		Attr1:                  req.Attr1,
		Attr2:                  req.Attr2,
		Attr3:                  req.Attr3,
		Attr4:                  req.Attr4,
		CreateBy:               baizeContext.GetUserName(c),
		UpdateBy:               baizeContext.GetUserName(c),
		CreateById:             baizeContext.GetUserId(c),
		UpdateById:             baizeContext.GetUserId(c),
	}

	ret := dao.db.Create(iqc)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return &baize.EmptyResponse{}, nil
}

// Update 修改来料检验单
func (dao *QcIqcDaoImpl) Update(c *gin.Context, req *qcIqcApi.QcIqcUpdateReq) (*baize.EmptyResponse, error) {
	// 先查询原记录
	iqc, err := dao.GetById(c, req.IqcId)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.IqcName != nil {
		iqc.IqcName = *req.IqcName
	}
	if req.TemplateId != nil {
		iqc.TemplateId = *req.TemplateId
	}
	if req.VendorId != nil {
		iqc.VendorId = *req.VendorId
	}
	if req.VendorCode != nil {
		iqc.VendorCode = *req.VendorCode
	}
	if req.VendorName != nil {
		iqc.VendorName = *req.VendorName
	}
	if req.ItemId != nil {
		iqc.ItemId = *req.ItemId
	}
	if req.QuantityRecived != nil {
		iqc.QuantityRecived = *req.QuantityRecived
	}
	if req.Status != nil {
		iqc.Status = req.Status
	}
	iqc.UpdateBy = baizeContext.GetUserName(c)
	iqc.UpdateById = baizeContext.GetUserId(c)

	ret := dao.db.Save(iqc)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return &baize.EmptyResponse{}, nil
}

// Delete 批量删除来料检验单
func (dao *QcIqcDaoImpl) Delete(c *gin.Context, iqcIds []int64) error {
	if len(iqcIds) == 0 {
		return nil
	}

	ret := dao.db.Where("id IN ?", iqcIds).Delete(&qcIqcModel.QcIqc{})
	if ret.Error != nil {
		return ret.Error
	}

	return nil
}

// GetById 根据ID查询来料检验单
func (dao *QcIqcDaoImpl) GetById(c *gin.Context, iqcId int64) (*qcIqcModel.QcIqc, error) {
	var iqc qcIqcModel.QcIqc
	ret := dao.db.Where("id = ?", iqcId).First(&iqc)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &iqc, nil
}
