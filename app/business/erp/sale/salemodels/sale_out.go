package salemodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// SaleOut ERP 销售出库
type SaleOut struct {
	ID                int64      `json:"id,string" gorm:"column:id"`
	No                string     `json:"no" gorm:"column:no"`
	Status            int32      `json:"status" gorm:"column:status"`
	CustomerID        int64      `json:"customerId" gorm:"column:customer_id"`
	AccountID         int64      `json:"accountId" gorm:"column:account_id"`
	SaleUserID        int64      `json:"saleUserId" gorm:"column:sale_user_id"`
	OutTime           *time.Time `json:"outTime" gorm:"column:out_time"`
	OrderID           int64      `json:"orderId" gorm:"column:order_id"`
	OrderNo           string     `json:"orderNo" gorm:"column:order_no"`
	TotalCount        float64    `json:"totalCount" gorm:"column:total_count"`
	TotalPrice        float64    `json:"totalPrice" gorm:"column:total_price"`
	ReceiptPrice      float64    `json:"receiptPrice" gorm:"column:receipt_price"`
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

// SaleOutUpsert ERP 销售出库新增修改参数
type SaleOutUpsert struct {
	ID                int64   `json:"id,string"`
	No                string  `json:"no" binding:"required" label:"销售出库单号"`
	Status            int32   `json:"status"`
	CustomerID        int64   `json:"customerId"`
	AccountID         int64   `json:"accountId"`
	SaleUserID        int64   `json:"saleUserId"`
	OutTime           string  `json:"outTime"`
	OrderID           int64   `json:"orderId"`
	OrderNo           string  `json:"orderNo"`
	TotalCount        float64 `json:"totalCount"`
	TotalPrice        float64 `json:"totalPrice"`
	ReceiptPrice      float64 `json:"receiptPrice"`
	TotalProductPrice float64 `json:"totalProductPrice"`
	TotalTaxPrice     float64 `json:"totalTaxPrice"`
	DiscountPercent   float64 `json:"discountPercent"`
	DiscountPrice     float64 `json:"discountPrice"`
	OtherPrice        float64 `json:"otherPrice"`
	FileURL           string  `json:"fileUrl"`
	Remark            string  `json:"remark"`
}

// SaleOutQuery ERP 销售出库查询参数
type SaleOutQuery struct {
	No         string `form:"no" filter:"like,no"`
	Status     *int32 `form:"status" filter:"eq,status"`
	CustomerID int64  `form:"customerId" filter:"eq,customer_id"`
	AccountID  int64  `form:"accountId" filter:"eq,account_id"`
	OrderID    int64  `form:"orderId" filter:"eq,order_id"`
	Page       int64  `form:"page"`
	Size       int64  `form:"size"`
}

// SaleOutListData ERP 销售出库分页数据
type SaleOutListData struct {
	Rows  []*SaleOut `json:"rows"`
	Total int64      `json:"total"`
}
