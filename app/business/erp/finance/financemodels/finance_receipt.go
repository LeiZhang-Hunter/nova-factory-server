package financemodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// FinanceReceipt ERP 收款单
type FinanceReceipt struct {
	ID            int64      `json:"id,string" gorm:"column:id"`
	No            string     `json:"no" gorm:"column:no"`
	Status        int32      `json:"status" gorm:"column:status"`
	ReceiptTime   *time.Time `json:"receiptTime" gorm:"column:receipt_time"`
	FinanceUserID int64      `json:"financeUserId" gorm:"column:finance_user_id"`
	CustomerID    int64      `json:"customerId" gorm:"column:customer_id"`
	AccountID     int64      `json:"accountId" gorm:"column:account_id"`
	TotalPrice    float64    `json:"totalPrice" gorm:"column:total_price"`
	DiscountPrice float64    `json:"discountPrice" gorm:"column:discount_price"`
	ReceiptPrice  float64    `json:"receiptPrice" gorm:"column:receipt_price"`
	Remark        string     `json:"remark" gorm:"column:remark"`
	DeptID        int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// FinanceReceiptUpsert ERP 收款单新增修改参数
type FinanceReceiptUpsert struct {
	ID            int64   `json:"id,string"`
	No            string  `json:"no" binding:"required" label:"收款单号"`
	Status        int32   `json:"status"`
	ReceiptTime   string  `json:"receiptTime"`
	FinanceUserID int64   `json:"financeUserId"`
	CustomerID    int64   `json:"customerId"`
	AccountID     int64   `json:"accountId"`
	TotalPrice    float64 `json:"totalPrice"`
	DiscountPrice float64 `json:"discountPrice"`
	ReceiptPrice  float64 `json:"receiptPrice"`
	Remark        string  `json:"remark"`
}

// FinanceReceiptQuery ERP 收款单查询参数
type FinanceReceiptQuery struct {
	No            string `form:"no" filter:"like,no"`
	Status        *int32 `form:"status" filter:"eq,status"`
	CustomerID    int64  `form:"customerId" filter:"eq,customer_id"`
	AccountID     int64  `form:"accountId" filter:"eq,account_id"`
	FinanceUserID int64  `form:"financeUserId" filter:"eq,finance_user_id"`
	Page          int64  `form:"page"`
	Size          int64  `form:"size"`
}

// FinanceReceiptListData ERP 收款单分页数据
type FinanceReceiptListData struct {
	Rows  []*FinanceReceipt `json:"rows"`
	Total int64             `json:"total"`
}
