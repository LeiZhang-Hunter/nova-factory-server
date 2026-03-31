package buildingdaoimpl

import (
	"encoding/json"
	"errors"
	"nova-factory-server/app/business/iot/asset/building/buildingdao"
	"nova-factory-server/app/business/iot/asset/building/buildingmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FloorDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewFloorDaoImpl(db *gorm.DB) buildingdao.FloorDao {
	return &FloorDaoImpl{
		db:    db,
		table: "sys_floor",
	}
}

func (b *FloorDaoImpl) Set(c *gin.Context, data *buildingmodels.SetSysFloor) (*buildingmodels.SysFloor, error) {

	value := buildingmodels.FromSetSysFloorToSysFloor(data)
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

func (b *FloorDaoImpl) List(c *gin.Context, req *buildingmodels.SetSysFloorListReq) (*buildingmodels.SetSysFloorList, error) {
	db := b.db.Table(b.table)

	if req != nil && req.FloorName != "" {
		db = db.Where("floor_name like ?", "%"+req.FloorName+"%")
	}
	if req != nil && req.BuildingID != 0 {
		db = db.Where("building_id = ?", req.BuildingID)
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
	var dto []*buildingmodels.SysFloor

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &buildingmodels.SetSysFloorList{
			Rows:  make([]*buildingmodels.SysFloor, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &buildingmodels.SetSysFloorList{
			Rows:  make([]*buildingmodels.SysFloor, 0),
			Total: 0,
		}, ret.Error
	}

	return &buildingmodels.SetSysFloorList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (b *FloorDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := b.db.Table(b.table).Where("id = ?", ids).Delete(&buildingmodels.SysFloor{})
	return ret.Error
}

func (b *FloorDaoImpl) GetByIds(c *gin.Context, ids []uint64) ([]*buildingmodels.SysFloor, error) {
	var list []*buildingmodels.SysFloor
	ret := b.db.Table(b.table).Where("id in (?)", ids).Find(&list)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return list, nil
}

func (b *FloorDaoImpl) CheckUniqueFloor(c *gin.Context, id int64, buildingId int64, level int8) (int64, error) {
	var count int64
	check := b.db.Table(b.table).Where("building_id = ? AND level = ? AND state = ?", buildingId, level, commonStatus.NORMAL)
	if id != 0 {
		check = check.Where("id != ?", id)
	}
	ret := check.Count(&count)
	return count, ret.Error
}

func (b *FloorDaoImpl) SaveLayout(c *gin.Context, id int64, layout *buildingmodels.FloorLayout) error {
	if layout == nil {
		return errors.New("layout == nil")
	}
	content, err := json.Marshal(layout)
	if err != nil {
		return err
	}
	ret := b.db.Table(b.table).Where("id = ?", id).Update("layout", content)
	return ret.Error
}

func (b *FloorDaoImpl) Info(c *gin.Context, id int64) (*buildingmodels.SysFloor, error) {
	var info buildingmodels.SysFloor
	ret := b.db.Table(b.table).Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if len(info.Layout) > 0 {
		var layout buildingmodels.FloorLayout
		err := json.Unmarshal([]byte(info.Layout), &layout)
		if err != nil {
			return nil, err
		}
		info.LayoutData = &layout
	} else {
		info.LayoutData = &buildingmodels.FloorLayout{
			FloorId: 0,
			Zones:   make([]buildingmodels.Zone, 0),
		}
	}

	return &info, nil
}

func (b *FloorDaoImpl) All(c *gin.Context) ([]*buildingmodels.SysFloor, error) {
	var list []*buildingmodels.SysFloor
	ret := b.db.Table(b.table).Find(&list)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return list, nil
		}
		return nil, ret.Error
	}
	return list, nil
}
