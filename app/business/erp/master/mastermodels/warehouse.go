package mastermodels

import (
	"nova-factory-server/app/baize"
)

// Warehouse ERP 仓库
type Warehouse struct {
	ID             int64   `json:"id,string" gorm:"column:id"`
	Name           string  `json:"name" gorm:"column:name"`
	Address        string  `json:"address" gorm:"column:address"`
	Sort           int64   `json:"sort" gorm:"column:sort"`
	Remark         string  `json:"remark" gorm:"column:remark"`
	Principal      string  `json:"principal" gorm:"column:principal"`
	WarehousePrice float64 `json:"warehousePrice" gorm:"column:warehouse_price"`
	TruckagePrice  float64 `json:"truckagePrice" gorm:"column:truckage_price"`
	Status         int32   `json:"status" gorm:"column:status"`
	DefaultStatus  *bool   `json:"defaultStatus" gorm:"column:default_status"`
	DeptID         int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// WarehouseUpsert ERP 仓库新增修改参数
type WarehouseUpsert struct {
	ID             int64   `json:"id,string"`
	Name           string  `json:"name" binding:"required" label:"仓库名称"`
	Address        string  `json:"address"`
	Sort           int64   `json:"sort"`
	Remark         string  `json:"remark"`
	Principal      string  `json:"principal"`
	WarehousePrice float64 `json:"warehousePrice"`
	TruckagePrice  float64 `json:"truckagePrice"`
	Status         int32   `json:"status"`
	DefaultStatus  *bool   `json:"defaultStatus"`
}

// WarehouseQuery ERP 仓库查询参数
type WarehouseQuery struct {
	Name          string `form:"name" filter:"like,name"`
	Status        *int32 `form:"status" filter:"eq,status"`
	DefaultStatus *bool  `form:"defaultStatus" filter:"eq,default_status"`
	Page          int64  `form:"page"`
	Size          int64  `form:"size"`
}

// WarehouseListData ERP 仓库分页数据
type WarehouseListData struct {
	Rows  []*Warehouse `json:"rows"`
	Total int64        `json:"total"`
}
