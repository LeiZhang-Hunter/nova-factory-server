package craftRouteDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type SysProRouteProductBomDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewSysProRouteProductBomDaoImpl(db *gorm.DB) craftRouteDao.ISysProRouteProductBomDao {
	return &SysProRouteProductBomDaoImpl{
		db:        db,
		tableName: "pro_route_product_bom",
	}
}

func (s *SysProRouteProductBomDaoImpl) Add(c *gin.Context, info *craftRouteModels.SysSetProRouteProductBom) (*craftRouteModels.SysProRouteProductBom, error) {
	data := craftRouteModels.NewSysProRouteProductBom(info)
	data.RecordID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	ret := s.db.Table(s.tableName).Create(data)
	return data, ret.Error
}

func (s *SysProRouteProductBomDaoImpl) Update(c *gin.Context, info *craftRouteModels.SysSetProRouteProductBom) (*craftRouteModels.SysProRouteProductBom, error) {
	data := craftRouteModels.NewSysProRouteProductBom(info)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	ret := s.db.Table(s.tableName).Where("record_id = ?", data.RecordID).Updates(data)
	return data, ret.Error
}

func (s *SysProRouteProductBomDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := s.db.Table(s.tableName).Where("record_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (s *SysProRouteProductBomDaoImpl) List(c *gin.Context, req *craftRouteModels.SysProRouteProductBomReq) (*craftRouteModels.SysProRouteProductBomList, error) {
	db := s.db.Table(s.tableName).Table(s.tableName)

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
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*craftRouteModels.SysProRouteProductBom

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProRouteProductBomList{
			Rows:  make([]*craftRouteModels.SysProRouteProductBom, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProRouteProductBomList{
			Rows:  make([]*craftRouteModels.SysProRouteProductBom, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysProRouteProductBomList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (s *SysProRouteProductBomDaoImpl) GetByRouteId(c *gin.Context, routeID int64) ([]*craftRouteModels.SysProRouteProductBom, error) {
	var boms []*craftRouteModels.SysProRouteProductBom
	ret := s.db.Table(s.tableName).Where("route_id = ?", routeID).Where("state = ?", commonStatus.NORMAL).Find(&boms)
	return boms, ret.Error
}
