package buildingDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/building/buildingModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type BuildingDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewBuildingDaoImpl(db *gorm.DB) buildingDao.BuildingDao {
	return &BuildingDaoImpl{
		db:    db,
		table: "sys_building",
	}
}

func (b *BuildingDaoImpl) Set(c *gin.Context, data *buildingModels.SetSysBuilding) (*buildingModels.SysBuilding, error) {
	value := buildingModels.FromSetSysBuildingToSysBuilding(data)
	if value.ID == 0 {
		value.SetCreateBy(baizeContext.GetUserId(c))
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		ret := b.db.Table(b.table).Create(&value)
		return value, ret.Error
	} else {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := b.db.Table(b.table).Debug().Where("id = ?", value.ID).Updates(&value)
		return value, ret.Error
	}
}

func (b *BuildingDaoImpl) List(c *gin.Context, req *buildingModels.SetSysBuildingListReq) (*buildingModels.SetSysBuildingList, error) {
	db := b.db.Table(b.table)

	if req != nil && req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
	}
	if req != nil && req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req != nil && req.Status != "" {
		db = db.Where("status = ?", req.Status)
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
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*buildingModels.SysBuilding

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &buildingModels.SetSysBuildingList{
			Rows:  make([]*buildingModels.SysBuilding, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &buildingModels.SetSysBuildingList{
			Rows:  make([]*buildingModels.SysBuilding, 0),
			Total: 0,
		}, ret.Error
	}

	return &buildingModels.SetSysBuildingList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (b *BuildingDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := b.db.Table(b.table).Where("id = ?", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
