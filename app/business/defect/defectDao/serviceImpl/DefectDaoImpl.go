package serviceImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/defect/defectApi"
	"nova-factory-server/app/business/defect/defectDao"
	"nova-factory-server/app/business/defect/defectModel"
	"nova-factory-server/app/utils/baizeContext"
)

type DefectDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewDefectDaoImpl(db *gorm.DB) defectDao.DefectDao {
	return &DefectDaoImpl{
		db:        db,
		tableName: "sys_defect",
	}
}

func (dao *DefectDaoImpl) List(c *gin.Context, req *defectApi.DefectListReq) (int64, []*defectModel.Defect, error) {
	db := dao.db.Table(dao.tableName)
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
	var dto []*defectModel.Defect
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

func (dao *DefectDaoImpl) Create(c *gin.Context, req *defectApi.DefectCreateReq) (*baize.EmptyResponse, error) {
	defect := &defectModel.Defect{
		DefectCode:  req.DefectCode,
		DefectName:  req.DefectName,
		IndexType:   req.IndexType,
		DefectLevel: req.DefectLevel,
		DeptId:      baizeContext.GetDeptId(c),
		State:       "0", // 默认状态
		CreateBy:    baizeContext.GetUserName(c),
		CreateById:  baizeContext.GetUserId(c),
		UpdateById:  baizeContext.GetUserId(c),
		UpdateBy:    baizeContext.GetUserName(c),
		Remark:      req.Remark,
		Attr1:       req.Attr1,
		Attr2:       req.Attr2,
		Attr3:       req.Attr3,
		Attr4:       req.Attr4,
	}

	ret := dao.db.Table(dao.tableName).Create(defect)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return &baize.EmptyResponse{}, nil
}

func (dao *DefectDaoImpl) Update(c *gin.Context, req *defectApi.DefectUpdateReq) (*baize.EmptyResponse, error) {
	// 先查询原记录
	_, err := dao.GetById(c, req.DefectId)
	if err != nil {
		return nil, err
	}
	// 更新字段
	updateData := map[string]interface{}{
		"defect_code":  req.DefectCode,
		"defect_name":  req.DefectName,
		"index_type":   req.IndexType,
		"defect_level": req.DefectLevel,
		"state":        req.State,
		"remark":       req.Remark,
		"attr_1":       req.Attr1,
		"attr_2":       req.Attr2,
		"attr_3":       req.Attr3,
		"attr_4":       req.Attr4,
		"update_by_id": baizeContext.GetUserId(c),
		"update_by":    baizeContext.GetUserName(c),
	}

	ret := dao.db.Table(dao.tableName).Where("id = ?", req.DefectId).Updates(updateData)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// 返回更新后的数据
	return &baize.EmptyResponse{}, nil
}

func (dao *DefectDaoImpl) Delete(c *gin.Context, defectIds []int64) error {
	if len(defectIds) == 0 {
		return nil
	}

	ret := dao.db.Table(dao.tableName).Where("id IN ?", defectIds).Delete(&defectApi.DefectUpdateRes{})
	if ret.Error != nil {
		return ret.Error
	}

	return nil
}

func (dao *DefectDaoImpl) GetById(c *gin.Context, defectId int64) (*defectModel.Defect, error) {
	var defect defectModel.Defect
	ret := dao.db.Table(dao.tableName).Where("id = ?", defectId).First(&defect)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &defect, nil
}
