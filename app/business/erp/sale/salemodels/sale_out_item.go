package salemodels

import (
	"nova-factory-server/app/baize"
)

// SaleOutItem ERP 销售出库项
type SaleOutItem struct {
	ID            int64   `json:"id,string" gorm:"column:id"`
	OutID         int64   `json:"outId" gorm:"column:out_id"`
	OrderItemID   int64   `json:"orderItemId" gorm:"column:order_item_id"`
	WarehouseID   int64   `json:"warehouseId" gorm:"column:warehouse_id"`
	ProductID     int64   `json:"productId" gorm:"column:product_id"`
	ProductUnitID int64   `json:"productUnitId" gorm:"column:product_unit_id"`
	ProductPrice  float64 `json:"productPrice" gorm:"column:product_price"`
	Count         float64 `json:"count" gorm:"column:count"`
	TotalPrice    float64 `json:"totalPrice" gorm:"column:total_price"`
	TaxPercent    float64 `json:"taxPercent" gorm:"column:tax_percent"`
	TaxPrice      float64 `json:"taxPrice" gorm:"column:tax_price"`
	Remark        string  `json:"remark" gorm:"column:remark"`
	DeptID        int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// SaleOutItemUpsert ERP 销售出库项新增修改参数
type SaleOutItemUpsert struct {
	ID            int64   `json:"id,string"`
	OutID         int64   `json:"outId"`
	OrderItemID   int64   `json:"orderItemId"`
	WarehouseID   int64   `json:"warehouseId"`
	ProductID     int64   `json:"productId"`
	ProductUnitID int64   `json:"productUnitId"`
	ProductPrice  float64 `json:"productPrice"`
	Count         float64 `json:"count"`
	TotalPrice    float64 `json:"totalPrice"`
	TaxPercent    float64 `json:"taxPercent"`
	TaxPrice      float64 `json:"taxPrice"`
	Remark        string  `json:"remark"`
}

// SaleOutItemQuery ERP 销售出库项查询参数
type SaleOutItemQuery struct {
	WarehouseID int64 `form:"warehouseId" filter:"eq,warehouse_id"`
	ProductID   int64 `form:"productId" filter:"eq,product_id"`
	OutID       int64 `form:"outId" filter:"eq,out_id"`
	Page        int64 `form:"page"`
	Size        int64 `form:"size"`
}

// SaleOutItemListData ERP 销售出库项分页数据
type SaleOutItemListData struct {
	Rows  []*SaleOutItem `json:"rows"`
	Total int64          `json:"total"`
}
