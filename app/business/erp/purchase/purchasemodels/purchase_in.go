package purchasemodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// PurchaseIn ERP 采购入库
type PurchaseIn struct {
	ID                int64      `json:"id,string" gorm:"column:id"`
	No                string     `json:"no" gorm:"column:no"`
	Status            int32      `json:"status" gorm:"column:status"`
	SupplierID        int64      `json:"supplierId" gorm:"column:supplier_id"`
	AccountID         int64      `json:"accountId" gorm:"column:account_id"`
	InTime            *time.Time `json:"inTime" gorm:"column:in_time"`
	OrderID           int64      `json:"orderId" gorm:"column:order_id"`
	OrderNo           string     `json:"orderNo" gorm:"column:order_no"`
	TotalCount        float64    `json:"totalCount" gorm:"column:total_count"`
	TotalPrice        float64    `json:"totalPrice" gorm:"column:total_price"`
	PaymentPrice      float64    `json:"paymentPrice" gorm:"column:payment_price"`
	TotalProductPrice float64    `json:"totalProductPrice" gorm:"column:total_product_price"`
	TotalTaxPrice     float64    `json:"totalTaxPrice" gorm:"column:total_tax_price"`
	DiscountPercent   float64    `json:"discountPercent" gorm:"column:discount_percent"`
	DiscountPrice     float64    `json:"discountPrice" gorm:"column:discount_price"`
	OtherPrice        float64    `json:"otherPrice" gorm:"column:other_price"`
	FileURL           string     `json:"fileUrl" gorm:"column:file_url"`
	Remark            string     `json:"remark" gorm:"column:remark"`
	DeptID            int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// PurchaseInUpsert ERP 采购入库新增修改参数
type PurchaseInUpsert struct {
	ID                int64   `json:"id,string"`
	No                string  `json:"no" binding:"required" label:"采购入库单号"`
	Status            int32   `json:"status"`
	SupplierID        int64   `json:"supplierId"`
	AccountID         int64   `json:"accountId"`
	InTime            string  `json:"inTime"`
	OrderID           int64   `json:"orderId"`
	OrderNo           string  `json:"orderNo"`
	TotalCount        float64 `json:"totalCount"`
	TotalPrice        float64 `json:"totalPrice"`
	PaymentPrice      float64 `json:"paymentPrice"`
	TotalProductPrice float64 `json:"totalProductPrice"`
	TotalTaxPrice     float64 `json:"totalTaxPrice"`
	DiscountPercent   float64 `json:"discountPercent"`
	DiscountPrice     float64 `json:"discountPrice"`
	OtherPrice        float64 `json:"otherPrice"`
	FileURL           string  `json:"fileUrl"`
	Remark            string  `json:"remark"`
}

// PurchaseInQuery ERP 采购入库查询参数
type PurchaseInQuery struct {
	No         string `form:"no" filter:"like,no"`
	Status     *int32 `form:"status" filter:"eq,status"`
	SupplierID int64  `form:"supplierId" filter:"eq,supplier_id"`
	AccountID  int64  `form:"accountId" filter:"eq,account_id"`
	OrderID    int64  `form:"orderId" filter:"eq,order_id"`
	Page       int64  `form:"page"`
	Size       int64  `form:"size"`
}

// PurchaseInListData ERP 采购入库分页数据
type PurchaseInListData struct {
	Rows  []*PurchaseIn `json:"rows"`
	Total int64         `json:"total"`
}
