package financemodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// FinancePayment ERP 付款单
type FinancePayment struct {
	ID            int64      `json:"id,string" gorm:"column:id"`
	No            string     `json:"no" gorm:"column:no"`
	Status        int32      `json:"status" gorm:"column:status"`
	PaymentTime   *time.Time `json:"paymentTime" gorm:"column:payment_time"`
	FinanceUserID int64      `json:"financeUserId" gorm:"column:finance_user_id"`
	SupplierID    int64      `json:"supplierId" gorm:"column:supplier_id"`
	AccountID     int64      `json:"accountId" gorm:"column:account_id"`
	TotalPrice    float64    `json:"totalPrice" gorm:"column:total_price"`
	DiscountPrice float64    `json:"discountPrice" gorm:"column:discount_price"`
	PaymentPrice  float64    `json:"paymentPrice" gorm:"column:payment_price"`
	Remark        string     `json:"remark" gorm:"column:remark"`
	DeptID        int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// FinancePaymentUpsert ERP 付款单新增修改参数
type FinancePaymentUpsert struct {
	ID            int64   `json:"id,string"`
	No            string  `json:"no" binding:"required" label:"付款单号"`
	Status        int32   `json:"status"`
	PaymentTime   string  `json:"paymentTime"`
	FinanceUserID int64   `json:"financeUserId"`
	SupplierID    int64   `json:"supplierId"`
	AccountID     int64   `json:"accountId"`
	TotalPrice    float64 `json:"totalPrice"`
	DiscountPrice float64 `json:"discountPrice"`
	PaymentPrice  float64 `json:"paymentPrice"`
	Remark        string  `json:"remark"`
}

// FinancePaymentQuery ERP 付款单查询参数
type FinancePaymentQuery struct {
	No            string `form:"no" filter:"like,no"`
	Status        *int32 `form:"status" filter:"eq,status"`
	SupplierID    int64  `form:"supplierId" filter:"eq,supplier_id"`
	AccountID     int64  `form:"accountId" filter:"eq,account_id"`
	FinanceUserID int64  `form:"financeUserId" filter:"eq,finance_user_id"`
	Page          int64  `form:"page"`
	Size          int64  `form:"size"`
}

// FinancePaymentListData ERP 付款单分页数据
type FinancePaymentListData struct {
	Rows  []*FinancePayment `json:"rows"`
	Total int64             `json:"total"`
}
