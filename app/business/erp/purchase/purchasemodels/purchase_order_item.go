package purchasemodels

import (
	"nova-factory-server/app/baize"
)

// PurchaseOrderItem ERP 采购订单项
type PurchaseOrderItem struct {
	ID            int64   `json:"id,string" gorm:"column:id"`
	OrderID       int64   `json:"orderId" gorm:"column:order_id"`
	ProductID     int64   `json:"productId" gorm:"column:product_id"`
	ProductUnitID int64   `json:"productUnitId" gorm:"column:product_unit_id"`
	ProductPrice  float64 `json:"productPrice" gorm:"column:product_price"`
	Count         float64 `json:"count" gorm:"column:count"`
	TotalPrice    float64 `json:"totalPrice" gorm:"column:total_price"`
	TaxPercent    float64 `json:"taxPercent" gorm:"column:tax_percent"`
	TaxPrice      float64 `json:"taxPrice" gorm:"column:tax_price"`
	Remark        string  `json:"remark" gorm:"column:remark"`
	InCount       float64 `json:"inCount" gorm:"column:in_count"`
	ReturnCount   float64 `json:"returnCount" gorm:"column:return_count"`
	DeptID        int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// PurchaseOrderItemUpsert ERP 采购订单项新增修改参数
type PurchaseOrderItemUpsert struct {
	ID            int64   `json:"id,string"`
	OrderID       int64   `json:"orderId"`
	ProductID     int64   `json:"productId" binding:"required" label:"产品"`
	ProductUnitID int64   `json:"productUnitId"`
	ProductPrice  float64 `json:"productPrice" binding:"required" label:"采购单价"`
	Count         float64 `json:"count" binding:"required" label:"采购数量"`
	TotalPrice    float64 `json:"totalPrice"`
	TaxPercent    float64 `json:"taxPercent"`
	TaxPrice      float64 `json:"taxPrice"`
	Remark        string  `json:"remark"`
	InCount       float64 `json:"inCount"`
	ReturnCount   float64 `json:"returnCount"`
}

// PurchaseOrderItemQuery ERP 采购订单项查询参数
type PurchaseOrderItemQuery struct {
	OrderID   int64 `form:"orderId" filter:"eq,order_id"`
	ProductID int64 `form:"productId" filter:"eq,product_id"`
	Page      int64 `form:"page"`
	Size      int64 `form:"size"`
}

// PurchaseOrderItemListData ERP 采购订单项分页数据
type PurchaseOrderItemListData struct {
	Rows  []*PurchaseOrderItem `json:"rows"`
	Total int64                `json:"total"`
}
