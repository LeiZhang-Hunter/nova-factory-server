package stockmodels

import (
	"nova-factory-server/app/baize"
)

// Stock ERP 产品库存
type Stock struct {
	ID          int64   `json:"id,string" gorm:"column:id"`
	ProductID   int64   `json:"productId" gorm:"column:product_id"`
	WarehouseID int64   `json:"warehouseId" gorm:"column:warehouse_id"`
	Count       float64 `json:"count" gorm:"column:count"`
	DeptID      int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockUpsert ERP 产品库存新增修改参数
type StockUpsert struct {
	ID          int64   `json:"id,string"`
	ProductID   int64   `json:"productId"`
	WarehouseID int64   `json:"warehouseId"`
	Count       float64 `json:"count"`
}

// StockQuery ERP 产品库存查询参数
type StockQuery struct {
	WarehouseID int64 `form:"warehouseId" filter:"eq,warehouse_id"`
	ProductID   int64 `form:"productId" filter:"eq,product_id"`
	Page        int64 `form:"page"`
	Size        int64 `form:"size"`
}

// StockListData ERP 产品库存分页数据
type StockListData struct {
	Rows  []*Stock `json:"rows"`
	Total int64    `json:"total"`
}
