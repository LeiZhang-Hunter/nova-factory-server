package stockmodels

import (
	"nova-factory-server/app/baize"
)

// StockCheckItem ERP 库存盘点单项
type StockCheckItem struct {
	ID            int64   `json:"id,string" gorm:"column:id"`
	CheckID       int64   `json:"checkId" gorm:"column:check_id"`
	WarehouseID   int64   `json:"warehouseId" gorm:"column:warehouse_id"`
	ProductID     int64   `json:"productId" gorm:"column:product_id"`
	ProductUnitID int64   `json:"productUnitId" gorm:"column:product_unit_id"`
	ProductPrice  float64 `json:"productPrice" gorm:"column:product_price"`
	StockCount    float64 `json:"stockCount" gorm:"column:stock_count"`
	ActualCount   float64 `json:"actualCount" gorm:"column:actual_count"`
	Count         float64 `json:"count" gorm:"column:count"`
	TotalPrice    float64 `json:"totalPrice" gorm:"column:total_price"`
	Remark        string  `json:"remark" gorm:"column:remark"`
	DeptID        int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockCheckItemUpsert ERP 库存盘点单项新增修改参数
type StockCheckItemUpsert struct {
	ID            int64   `json:"id,string"`
	CheckID       int64   `json:"checkId"`
	WarehouseID   int64   `json:"warehouseId"`
	ProductID     int64   `json:"productId"`
	ProductUnitID int64   `json:"productUnitId"`
	ProductPrice  float64 `json:"productPrice"`
	StockCount    float64 `json:"stockCount"`
	ActualCount   float64 `json:"actualCount"`
	Count         float64 `json:"count"`
	TotalPrice    float64 `json:"totalPrice"`
	Remark        string  `json:"remark"`
}

// StockCheckItemQuery ERP 库存盘点单项查询参数
type StockCheckItemQuery struct {
	WarehouseID int64 `form:"warehouseId" filter:"eq,warehouse_id"`
	ProductID   int64 `form:"productId" filter:"eq,product_id"`
	CheckID     int64 `form:"checkId" filter:"eq,check_id"`
	Page        int64 `form:"page"`
	Size        int64 `form:"size"`
}

// StockCheckItemListData ERP 库存盘点单项分页数据
type StockCheckItemListData struct {
	Rows  []*StockCheckItem `json:"rows"`
	Total int64             `json:"total"`
}
