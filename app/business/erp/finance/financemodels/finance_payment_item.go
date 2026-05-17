package financemodels

import (
	"nova-factory-server/app/baize"
)

// FinancePaymentItem ERP 付款项
type FinancePaymentItem struct {
	ID           int64   `json:"id,string" gorm:"column:id"`
	PaymentID    int64   `json:"paymentId" gorm:"column:payment_id"`
	BizType      int32   `json:"bizType" gorm:"column:biz_type"`
	BizID        int64   `json:"bizId" gorm:"column:biz_id"`
	BizNo        string  `json:"bizNo" gorm:"column:biz_no"`
	TotalPrice   float64 `json:"totalPrice" gorm:"column:total_price"`
	PaidPrice    float64 `json:"paidPrice" gorm:"column:paid_price"`
	PaymentPrice float64 `json:"paymentPrice" gorm:"column:payment_price"`
	Remark       string  `json:"remark" gorm:"column:remark"`
	DeptID       int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// FinancePaymentItemUpsert ERP 付款项新增修改参数
type FinancePaymentItemUpsert struct {
	ID           int64   `json:"id,string"`
	PaymentID    int64   `json:"paymentId"`
	BizType      int32   `json:"bizType"`
	BizID        int64   `json:"bizId"`
	BizNo        string  `json:"bizNo"`
	TotalPrice   float64 `json:"totalPrice"`
	PaidPrice    float64 `json:"paidPrice"`
	PaymentPrice float64 `json:"paymentPrice"`
	Remark       string  `json:"remark"`
}

// FinancePaymentItemQuery ERP 付款项查询参数
type FinancePaymentItemQuery struct {
	PaymentID int64  `form:"paymentId" filter:"eq,payment_id"`
	BizType   *int32 `form:"bizType" filter:"eq,biz_type"`
	Page      int64  `form:"page"`
	Size      int64  `form:"size"`
}

// FinancePaymentItemListData ERP 付款项分页数据
type FinancePaymentItemListData struct {
	Rows  []*FinancePaymentItem `json:"rows"`
	Total int64                 `json:"total"`
}
