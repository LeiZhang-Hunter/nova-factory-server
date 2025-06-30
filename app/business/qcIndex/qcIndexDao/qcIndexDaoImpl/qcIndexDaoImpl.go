package qcIndexDaoImpl

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIndex/qcIndexApi"
	"nova-factory-server/app/business/qcIndex/qcIndexDao"
	"nova-factory-server/app/business/qcIndex/qcIndexModel"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QcIndexDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewQcIndexDaoImpl(db *gorm.DB) qcIndexDao.QcIndexDao {
	return &QcIndexDaoImpl{
		db:        db,
		tableName: "qc_index",
	}
}

func (dao *QcIndexDaoImpl) List(c *gin.Context, req *qcIndexApi.QcIndexListReq) (int64, []*qcIndexModel.QcIndex, error) {
	db := dao.db.Model(qcIndexModel.QcIndex{})

	// 添加查询条件
	if req.IndexCode != "" {
		db = db.Where("index_code LIKE ?", "%"+req.IndexCode+"%")
	}
	if req.IndexName != "" {
		db = db.Where("index_name LIKE ?", "%"+req.IndexName+"%")
	}
	if req.IndexType != "" {
		db = db.Where("index_type = ?", req.IndexType)
	}
	if req.QcResultType != "" {
		db = db.Where("qc_result_type = ?", req.QcResultType)
	}

	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}

	var dto []*qcIndexModel.QcIndex
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

func (dao *QcIndexDaoImpl) Create(c *gin.Context, req *qcIndexApi.QcIndexCreateReq) (*baize.EmptyResponse, error) {
	index := &qcIndexModel.QcIndex{
		IndexCode:    req.IndexCode,
		IndexName:    req.IndexName,
		IndexType:    req.IndexType,
		QcTool:       req.QcTool,
		QcResultType: req.QcResultType,
		QcResultSpc:  req.QcResultSpc,
		Remark:       req.Remark,
		Attr1:        req.Attr1,
		Attr2:        req.Attr2,
		Attr3:        req.Attr3,
		Attr4:        req.Attr4,
		CreateBy:     baizeContext.GetUserName(c),
		UpdateBy:     baizeContext.GetUserName(c),
		CreateById:   baizeContext.GetUserId(c),
		UpdateById:   baizeContext.GetUserId(c),
	}

	ret := dao.db.Create(index)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return &baize.EmptyResponse{}, nil
}

func (dao *QcIndexDaoImpl) Update(c *gin.Context, req *qcIndexApi.QcIndexUpdateReq) (*baize.EmptyResponse, error) {
	// 先查询原记录
	index, err := dao.GetById(c, req.IndexId)
	if err != nil {
		return nil, err
	}

	// 更新字段
	index.IndexCode = req.IndexCode
	index.IndexName = req.IndexName
	index.IndexType = req.IndexType
	index.QcTool = req.QcTool
	index.QcResultType = req.QcResultType
	index.QcResultSpc = req.QcResultSpc
	index.Remark = req.Remark
	index.Attr1 = req.Attr1
	index.Attr2 = req.Attr2
	index.Attr3 = req.Attr3
	index.Attr4 = req.Attr4
	index.UpdateBy = baizeContext.GetUserName(c)
	index.UpdateById = baizeContext.GetUserId(c)

	ret := dao.db.Save(index)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return &baize.EmptyResponse{}, nil
}

func (dao *QcIndexDaoImpl) Delete(c *gin.Context, indexIds []int64) error {
	if len(indexIds) == 0 {
		return nil
	}

	ret := dao.db.Table(dao.tableName).Where("id IN ?", indexIds).Delete(&qcIndexModel.QcIndex{})
	if ret.Error != nil {
		return ret.Error
	}

	return nil
}

func (dao *QcIndexDaoImpl) GetById(c *gin.Context, indexId int64) (*qcIndexModel.QcIndex, error) {
	var index qcIndexModel.QcIndex
	ret := dao.db.Table(dao.tableName).Where("id = ?", indexId).First(&index)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &index, nil
}
