package stockmodels

import (
	"nova-factory-server/app/baize"
)

// StockMoveItem ERP 库存调拨单项
type StockMoveItem struct {
	ID              int64   `json:"id,string" gorm:"column:id"`
	MoveID          int64   `json:"moveId" gorm:"column:move_id"`
	FromWarehouseID int64   `json:"fromWarehouseID" gorm:"column:from_warehouse_id"`
	ToWarehouseID   int64   `json:"toWarehouseID" gorm:"column:to_warehouse_id"`
	ProductID       int64   `json:"productId" gorm:"column:product_id"`
	ProductUnitID   int64   `json:"productUnitId" gorm:"column:product_unit_id"`
	ProductPrice    float64 `json:"productPrice" gorm:"column:product_price"`
	Count           float64 `json:"count" gorm:"column:count"`
	TotalPrice      float64 `json:"totalPrice" gorm:"column:total_price"`
	Remark          string  `json:"remark" gorm:"column:remark"`
	DeptID          int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockMoveItemUpsert ERP 库存调拨单项新增修改参数
type StockMoveItemUpsert struct {
	ID              int64   `json:"id,string"`
	MoveID          int64   `json:"moveId"`
	FromWarehouseID int64   `json:"fromWarehouseID"`
	ToWarehouseID   int64   `json:"toWarehouseID"`
	ProductID       int64   `json:"productId"`
	ProductUnitID   int64   `json:"productUnitId"`
	ProductPrice    float64 `json:"productPrice"`
	Count           float64 `json:"count"`
	TotalPrice      float64 `json:"totalPrice"`
	Remark          string  `json:"remark"`
}

// StockMoveItemQuery ERP 库存调拨单项查询参数
type StockMoveItemQuery struct {
	ProductID int64 `form:"productId" filter:"eq,product_id"`
	MoveID    int64 `form:"moveId" filter:"eq,move_id"`
	Page      int64 `form:"page"`
	Size      int64 `form:"size"`
}

// StockMoveItemListData ERP 库存调拨单项分页数据
type StockMoveItemListData struct {
	Rows  []*StockMoveItem `json:"rows"`
	Total int64            `json:"total"`
}
