package models

import (
	"nova-factory-server/app/baize"
	"time"
)

// WarehouseArea WMS 库区。
type WarehouseArea struct {
	ID                    int64      `json:"id,string" gorm:"column:id"`
	WarehouseID           int64      `json:"warehouseId,string" gorm:"column:warehouse_id"`
	WarehouseName         string     `json:"warehouseName" gorm:"-"`
	AreaName              string     `json:"areaName" gorm:"column:area_name"`
	ParentID              int64      `json:"parentId,string" gorm:"column:parent_id"`
	ParentAreaName        string     `json:"parentAreaName" gorm:"-"`
	WarehouseAreaProperty int8       `json:"warehouseAreaProperty" gorm:"column:warehouse_area_property"`
	LastUpdateTime        *time.Time `json:"lastUpdateTime" gorm:"column:last_update_time"`
	IsValid               int8       `json:"isValid" gorm:"column:is_valid"`
	DeptID                int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// WarehouseAreaSet WMS 库区新增或修改参数。
type WarehouseAreaSet struct {
	ID                    int64  `json:"id,string"`
	WarehouseID           int64  `json:"warehouseId,string"`
	AreaName              string `json:"areaName"`
	ParentID              int64  `json:"parentId,string"`
	WarehouseAreaProperty int8   `json:"warehouseAreaProperty"`
	IsValid               int8   `json:"isValid"`
}

// WarehouseAreaQuery WMS 库区查询参数。
type WarehouseAreaQuery struct {
	WarehouseID           int64  `form:"warehouseId,string"`
	AreaName              string `form:"areaName"`
	ParentID              int64  `form:"parentId,string"`
	WarehouseAreaProperty *int8  `form:"warehouseAreaProperty"`
	IsValid               *int8  `form:"isValid"`
	Page                  int64  `form:"page"`
	Size                  int64  `form:"size"`
}

// WarehouseAreaListData WMS 库区分页数据。
type WarehouseAreaListData struct {
	Rows  []*WarehouseArea `json:"rows"`
	Total int64            `json:"total"`
}
