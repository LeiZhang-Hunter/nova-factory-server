package materialModels

import (
	"github.com/gogf/gf/util/gconv"
	"nova-factory-server/app/baize"
)

type MaterialVO struct {
	MaterialId   uint64  `json:"materialId,string" db:"material_id"`
	Name         string  `json:"name" db:"name" binding:"required"`
	Code         string  `json:"code" db:"code" binding:"required"`
	Model        string  `json:"model" db:"model" binding:"required"`
	Unit         string  `json:"unit" db:"device_building_id" binding:"required"`
	Factory      string  `json:"factory" db:"device_building_id" binding:"required"`
	Address      string  `json:"address" db:"device_building_id" binding:"required"`
	Price        float64 `json:"price" db:"device_building_id" binding:"required"`
	Total        uint64  `json:"total,string" db:"device_building_id"`
	Outbound     uint64  `json:"outbound,string" db:"device_building_id"`
	CurrentTotal uint64  `json:"currentTotal,string" db:"device_building_id"`
	DeptId       uint64  `json:"deptId,string" db:"dept_id"`
	State        int     `json:"ControlType" db:"control_type"`
	baize.BaseEntity
}

func NewMaterialVO(device *MaterialInfo) *MaterialVO {
	var vo MaterialVO
	gconv.Scan(device, &vo)
	return &vo
}

type MaterialValue struct {
	MaterialId     uint64  `json:"materialId,string" db:"material_id"`
	Name           string  `json:"name" db:"name" binding:"required"`
	Code           string  `json:"code" db:"code" binding:"required"`
	Model          string  `json:"model" db:"model" binding:"required"`
	Unit           string  `json:"unit" db:"device_building_id" binding:"required"`
	Factory        string  `json:"factory" db:"device_building_id" binding:"required"`
	Address        string  `json:"address" db:"device_building_id" binding:"required"`
	Price          float64 `json:"price" db:"device_building_id" binding:"required"`
	Total          uint64  `json:"total" db:"device_building_id"`
	Outbound       uint64  `json:"outbound" db:"device_building_id"`
	CurrentTotal   uint64  `json:"currentTotal" db:"device_building_id"`
	DeptId         uint64  `json:"deptId,string" db:"dept_id"`
	State          int     `json:"ControlType" db:"control_type"`
	CreateUserName string  `json:"createUserName"`
	UpdateUserName string  `json:"updateUserName"`
	baize.BaseEntity
}

type MaterialInfo struct {
	MaterialId uint64  `json:"materialId,string" form:"material_id" db:"material_id"`
	Name       string  `json:"name" db:"name" form:"name" binding:"required"`
	Code       string  `json:"code" db:"code" form:"code" binding:"required"`
	Model      string  `json:"model" db:"model" form:"model" binding:"required"`
	Unit       string  `json:"unit" db:"unit" form:"unit" binding:"required"`
	Factory    string  `json:"factory" db:"factory" form:"factory" binding:"required"`
	Address    string  `json:"address" db:"address" form:"address" binding:"required"`
	Price      float64 `json:"price"  db:"price" form:"price"`
}

type MaterialListReq struct {
	Name    string `json:"name,string" db:"name" form:"name"`
	Code    string `json:"code,string" db:"code" form:"code"`
	Model   string `json:"model,string" db:"model" form:"model"`
	Unit    string `json:"unit,string" db:"unit" form:"unit"`
	Factory string `json:"factory,string" db:"device_building_id" form:"factory"`
	Address string `json:"address,string" db:"device_building_id" form:"address"`
	baize.BaseEntityDQL
}

type MaterialInfoListData struct {
	Rows  []*MaterialVO `json:"rows"`
	Total int64         `json:"total"`
}

type MaterialInfoListValue struct {
	Rows  []*MaterialValue `json:"rows"`
	Total int64            `json:"total"`
}
