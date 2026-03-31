package buildingModels

import (
	"nova-factory-server/app/baize"
)

// SysFloor 楼层管理
type SysFloor struct {
	ID         int64        `gorm:"column:id;primaryKey;comment:楼层id" json:"id,string"`                   // 楼层id
	BuildingID int64        `gorm:"column:building_id;default:0;comment:建筑物id" json:"building_id,string"` // 建筑物id
	FloorName  string       `gorm:"column:floor_name;not null;comment:名称" json:"floor_name"`              // 名称
	Level      int8         `gorm:"column:level;not null;default:1;comment:楼层" json:"level"`              // 楼层
	DeptID     int64        `gorm:"column:dept_id;default:null;comment:部门ID" json:"dept_id"`              // 部门ID
	Layout     string       `gorm:"column:layout;comment:布局" json:"-"`                                    // 布局
	LayoutData *FloorLayout `gorm:"-" json:"layout"`                                                      // 布局
	baize.BaseEntity
	State int8 `gorm:"column:state;default:0;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func FromSetSysFloorToSysFloor(floor *SetSysFloor) *SysFloor {
	return &SysFloor{
		ID:         floor.ID,
		BuildingID: floor.BuildingID,
		FloorName:  floor.FloorName,
		Level:      floor.Level,
	}
}

type SetSysFloor struct {
	ID         int64  `gorm:"column:id;primaryKey;comment:楼层id" json:"id,string"`                   // 楼层id
	BuildingID int64  `gorm:"column:building_id;default:0;comment:建筑物id" json:"building_id,string"` // 建筑物id
	FloorName  string `gorm:"column:floor_name;not null;comment:名称" json:"floor_name"`              // 名称
	Level      int8   `gorm:"column:level;not null;default:1;comment:楼层" json:"level"`              // 楼层
}

type SetSysFloorListReq struct {
	FloorName  string `form:"floor_name"`
	BuildingID int64  `form:"building_id"`
	baize.BaseEntityDQL
}

type SetSysFloorList struct {
	Rows  []*SysFloor `json:"rows"`
	Total int64       `json:"total"`
}
type Zone struct {
	Name    string `json:"name"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Devices []struct {
		DeviceId   int64  `json:"deviceId,string"`
		DeviceName string `json:"deviceName"`
		X          int    `json:"x"`
		Y          int    `json:"y"`
	} `json:"devices"`
}
type FloorLayout struct {
	FloorId int64  `json:"floorId,string"`
	Zones   []Zone `json:"zones"`
}

type FloorLayoutInfoRequest struct {
	FloorId int64 `form:"floorId,string"`
}
