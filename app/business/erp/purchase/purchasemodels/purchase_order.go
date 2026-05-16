package purchasemodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// PurchaseOrder ERP 采购订单
type PurchaseOrder struct {
	ID                int64                `json:"id,string" gorm:"column:id"`
	No                string               `json:"no" gorm:"column:no"`
	Status            int32                `json:"status" gorm:"column:status"`
	SupplierID        int64                `json:"supplierId" gorm:"column:supplier_id"`
	AccountID         int64                `json:"accountId" gorm:"column:account_id"`
	OrderTime         *time.Time           `json:"orderTime" gorm:"column:order_time"`
	TotalCount        float64              `json:"totalCount" gorm:"column:total_count"`
	TotalPrice        float64              `json:"totalPrice" gorm:"column:total_price"`
	TotalProductPrice float64              `json:"totalProductPrice" gorm:"column:total_product_price"`
	TotalTaxPrice     float64              `json:"totalTaxPrice" gorm:"column:total_tax_price"`
	DiscountPercent   float64              `json:"discountPercent" gorm:"column:discount_percent"`
	DiscountPrice     float64              `json:"discountPrice" gorm:"column:discount_price"`
	DepositPrice      float64              `json:"depositPrice" gorm:"column:deposit_price"`
	FileURL           string               `json:"fileUrl" gorm:"column:file_url"`
	Remark            string               `json:"remark" gorm:"column:remark"`
	InCount           float64              `json:"inCount" gorm:"column:in_count"`
	ReturnCount       float64              `json:"returnCount" gorm:"column:return_count"`
	DeptID            int64                `json:"deptId" gorm:"column:dept_id"`
	Items             []*PurchaseOrderItem `json:"items,omitempty" gorm:"-"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// PurchaseOrderUpsert ERP 采购订单新增修改参数
type PurchaseOrderUpsert struct {
	ID              int64                      `json:"id,string"`
	No              string                     `json:"no"`
	SupplierID      int64                      `json:"supplierId" binding:"required" label:"供应商"`
	AccountID       int64                      `json:"accountId"`
	OrderTime       string                     `json:"orderTime" binding:"required" label:"下单时间"`
	DiscountPercent float64                    `json:"discountPercent"`
	DepositPrice    float64                    `json:"depositPrice"`
	FileURL         string                     `json:"fileUrl"`
	Remark          string                     `json:"remark"`
	Items           []*PurchaseOrderItemUpsert `json:"items" binding:"required" label:"采购订单明细"`
}

// PurchaseOrderStatusReq ERP 采购订单状态变更参数
type PurchaseOrderStatusReq struct {
	ID     int64 `json:"id,string" binding:"required" label:"采购订单ID"`
	Status int32 `json:"status" binding:"required" label:"状态"`
}

// PurchaseOrderQuery ERP 采购订单查询参数
type PurchaseOrderQuery struct {
	No         string `form:"no" filter:"like,no"`
	Status     *int32 `form:"status" filter:"eq,status"`
	SupplierID int64  `form:"supplierId" filter:"eq,supplier_id"`
	AccountID  int64  `form:"accountId" filter:"eq,account_id"`
	Page       int64  `form:"page"`
	Size       int64  `form:"size"`
}

// PurchaseOrderListData ERP 采购订单分页数据
type PurchaseOrderListData struct {
	Rows  []*PurchaseOrder `json:"rows"`
	Total int64            `json:"total"`
}
