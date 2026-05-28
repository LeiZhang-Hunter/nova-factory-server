package models

import (
	"nova-factory-server/app/baize"
	"time"
)

// WarehouseLocation WMS 库位。
type WarehouseLocation struct {
	ID                    int64      `json:"id,string" gorm:"column:id"`
	WarehouseID           int64      `json:"warehouseId,string" gorm:"column:warehouse_id"`
	WarehouseName         string     `json:"warehouseName" gorm:"column:warehouse_name"`
	WarehouseAreaID       int64      `json:"warehouseAreaId,string" gorm:"column:warehouse_area_id"`
	WarehouseAreaName     string     `json:"warehouseAreaName" gorm:"column:warehouse_area_name"`
	WarehouseAreaProperty int8       `json:"warehouseAreaProperty" gorm:"column:warehouse_area_property"`
	LocationName          string     `json:"locationName" gorm:"column:location_name"`
	LocationLength        float64    `json:"locationLength" gorm:"column:location_length"`
	LocationWidth         float64    `json:"locationWidth" gorm:"column:location_width"`
	LocationHeight        float64    `json:"locationHeight" gorm:"column:location_height"`
	LocationVolume        float64    `json:"locationVolume" gorm:"column:location_volume"`
	LocationLoad          float64    `json:"locationLoad" gorm:"column:location_load"`
	RoadwayNumber         string     `json:"roadwayNumber" gorm:"column:roadway_number"`
	ShelfNumber           string     `json:"shelfNumber" gorm:"column:shelf_number"`
	LayerNumber           string     `json:"layerNumber" gorm:"column:layer_number"`
	TagNumber             string     `json:"tagNumber" gorm:"column:tag_number"`
	LastUpdateTime        *time.Time `json:"lastUpdateTime" gorm:"column:last_update_time"`
	IsValid               int8       `json:"isValid" gorm:"column:is_valid"`
	DeptID                int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// WarehouseLocationSet WMS 库位新增或修改参数。
type WarehouseLocationSet struct {
	ID                    int64   `json:"id,string"`
	WarehouseID           int64   `json:"warehouseId,string"`
	WarehouseName         string  `json:"warehouseName"`
	WarehouseAreaID       int64   `json:"warehouseAreaId,string"`
	WarehouseAreaName     string  `json:"warehouseAreaName"`
	WarehouseAreaProperty int8    `json:"warehouseAreaProperty"`
	LocationName          string  `json:"locationName"`
	LocationLength        float64 `json:"locationLength"`
	LocationWidth         float64 `json:"locationWidth"`
	LocationHeight        float64 `json:"locationHeight"`
	LocationVolume        float64 `json:"locationVolume"`
	LocationLoad          float64 `json:"locationLoad"`
	RoadwayNumber         string  `json:"roadwayNumber"`
	ShelfNumber           string  `json:"shelfNumber"`
	LayerNumber           string  `json:"layerNumber"`
	TagNumber             string  `json:"tagNumber"`
	IsValid               int8    `json:"isValid"`
}

// WarehouseLocationQuery WMS 库位查询参数。
type WarehouseLocationQuery struct {
	WarehouseID     int64  `form:"warehouseId,string"`
	WarehouseAreaID int64  `form:"warehouseAreaId,string"`
	LocationName    string `form:"locationName"`
	RoadwayNumber   string `form:"roadwayNumber"`
	ShelfNumber     string `form:"shelfNumber"`
	LayerNumber     string `form:"layerNumber"`
	TagNumber       string `form:"tagNumber"`
	IsValid         *int8  `form:"isValid"`
	Page            int64  `form:"page"`
	Size            int64  `form:"size"`
}

// WarehouseLocationListData WMS 库位分页数据。
type WarehouseLocationListData struct {
	Rows  []*WarehouseLocation `json:"rows"`
	Total int64                `json:"total"`
}
