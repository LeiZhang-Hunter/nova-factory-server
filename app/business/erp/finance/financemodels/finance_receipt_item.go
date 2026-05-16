package financemodels

import (
	"nova-factory-server/app/baize"
)

// FinanceReceiptItem ERP 收款项
type FinanceReceiptItem struct {
	ID             int64   `json:"id,string" gorm:"column:id"`
	ReceiptID      int64   `json:"receiptId" gorm:"column:receipt_id"`
	BizType        int32   `json:"bizType" gorm:"column:biz_type"`
	BizID          int64   `json:"bizId" gorm:"column:biz_id"`
	BizNo          string  `json:"bizNo" gorm:"column:biz_no"`
	TotalPrice     float64 `json:"totalPrice" gorm:"column:total_price"`
	ReceiptedPrice float64 `json:"receiptedPrice" gorm:"column:receipted_price"`
	ReceiptPrice   float64 `json:"receiptPrice" gorm:"column:receipt_price"`
	Remark         string  `json:"remark" gorm:"column:remark"`
	DeptID         int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// FinanceReceiptItemUpsert ERP 收款项新增修改参数
type FinanceReceiptItemUpsert struct {
	ID             int64   `json:"id,string"`
	ReceiptID      int64   `json:"receiptId"`
	BizType        int32   `json:"bizType"`
	BizID          int64   `json:"bizId"`
	BizNo          string  `json:"bizNo"`
	TotalPrice     float64 `json:"totalPrice"`
	ReceiptedPrice float64 `json:"receiptedPrice"`
	ReceiptPrice   float64 `json:"receiptPrice"`
	Remark         string  `json:"remark"`
}

// FinanceReceiptItemQuery ERP 收款项查询参数
type FinanceReceiptItemQuery struct {
	ReceiptID int64  `form:"receiptId" filter:"eq,receipt_id"`
	BizType   *int32 `form:"bizType" filter:"eq,biz_type"`
	Page      int64  `form:"page"`
	Size      int64  `form:"size"`
}

// FinanceReceiptItemListData ERP 收款项分页数据
type FinanceReceiptItemListData struct {
	Rows  []*FinanceReceiptItem `json:"rows"`
	Total int64                 `json:"total"`
}
