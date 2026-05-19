package stockmodels

import (
	"nova-factory-server/app/baize"
)

// Stock ERP 产品库存
type Stock struct {
	ID            int64   `json:"id,string" gorm:"column:id"`
	ProductID     int64   `json:"productId,string" gorm:"column:product_id"`
	ProductName   string  `json:"productName" gorm:"-"`
	WarehouseID   int64   `json:"warehouseId,string" gorm:"column:warehouse_id"`
	WarehouseName string  `json:"warehouseName" gorm:"-"`
	Count         float64 `json:"count" gorm:"column:count"`
	EnableNum     float64 `json:"enablenum" gorm:"column:enablenum"`
	EnSaleNum     float64 `json:"ensalenum" gorm:"column:ensalenum"`
	DeptID        int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockUpsert ERP 产品库存新增修改参数
type StockUpsert struct {
	ID          int64   `json:"id,string"`
	ProductID   int64   `json:"productId,string"`
	WarehouseID int64   `json:"warehouseId,string"`
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
