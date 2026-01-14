package buildingModels

import (
	"nova-factory-server/app/baize"
)

// SysBuilding 设备管理
type SysBuilding struct {
	ID        int64       `gorm:"column:id;primaryKey;comment:id" json:"id,string"`        // id
	Name      string      `gorm:"column:name;comment:建筑物名称" json:"name"`                   // 建筑物名称
	Code      string      `gorm:"column:code;comment:建筑物编号" json:"code"`                   // 建筑物编号
	Type      string      `gorm:"column:type;comment:建筑物类型" json:"type"`                   // 建筑物类型
	Status    string      `gorm:"column:status;comment:建筑物状态" json:"status"`               // 建筑物状态
	Area      int64       `gorm:"column:area;not null;comment:建筑面积" json:"area"`           // 建筑面积
	Floors    int64       `gorm:"column:floors;not null;comment:楼层数" json:"floors"`        // 楼层数
	Year      string      `gorm:"column:year;comment:建造年份" json:"year"`                    // 建造年份
	LifeYears int64       `gorm:"column:lifeYears;not null;comment:使用年限" json:"lifeYears"` // 使用年限
	Manager   string      `gorm:"column:manager;comment:负责人w" json:"manager"`              // 负责人w
	Phone     string      `gorm:"column:phone;comment:单位" json:"phone"`                    // 单位
	Address   string      `gorm:"column:address;comment:单位" json:"address"`                // 单位
	Remark    string      `gorm:"column:remark;comment:单位" json:"remark"`                  // 单位
	DeptID    int64       `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`              // 部门ID
	Floor     []*SysFloor `gorm:"-" json:"floor"`
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func FromSetSysBuildingToSysBuilding(building *SetSysBuilding) *SysBuilding {
	return &SysBuilding{
		ID:        building.ID,
		Name:      building.Name,
		Code:      building.Code,
		Type:      building.Type,
		Status:    building.Status,
		Area:      building.Area,
		Floors:    building.Floors,
		Year:      building.Year,
		LifeYears: building.LifeYears,
		Manager:   building.Manager,
		Phone:     building.Phone,
		Address:   building.Address,
		Remark:    building.Remark,
	}
}

type SetSysBuilding struct {
	ID        int64  `gorm:"column:id;primaryKey;comment:id" json:"id,string"`        // id
	Name      string `gorm:"column:name;comment:建筑物名称" json:"name"`                   // 建筑物名称
	Code      string `gorm:"column:code;comment:建筑物编号" json:"code"`                   // 建筑物编号
	Type      string `gorm:"column:type;comment:建筑物类型" json:"type"`                   // 建筑物类型
	Status    string `gorm:"column:status;comment:建筑物状态" json:"status"`               // 建筑物状态
	Area      int64  `gorm:"column:area;not null;comment:建筑面积" json:"area"`           // 建筑面积
	Floors    int64  `gorm:"column:floors;not null;comment:楼层数" json:"floors"`        // 楼层数
	Year      string `gorm:"column:year;comment:建造年份" json:"year"`                    // 建造年份
	LifeYears int64  `gorm:"column:lifeYears;not null;comment:使用年限" json:"lifeYears"` // 使用年限
	Manager   string `gorm:"column:manager;comment:负责人w" json:"manager"`              // 负责人w
	Phone     string `gorm:"column:phone;comment:单位" json:"phone"`                    // 单位
	Address   string `gorm:"column:address;comment:单位" json:"address"`                // 单位
	Remark    string `gorm:"column:remark;comment:单位" json:"remark"`                  // 单位
}

type SetSysBuildingListReq struct {
	Name   string `form:"name"`
	Type   string `form:"type"`
	Status string `form:"status"`
	baize.BaseEntityDQL
}

type SetSysBuildingList struct {
	Rows  []*SysBuilding `json:"rows"`
	Total int64          `json:"total"`
}
