package models

import (
	"nova-factory-server/app/baize"
)

// OrderAccount ERP订单账户
type OrderAccount struct {
	ID          uint64  `json:"id,string" gorm:"column:id"`
	OrderID     uint64  `json:"order_id,string" gorm:"column:order_id"`
	Tid         string  `json:"tid" gorm:"column:tid"`
	FinanceCode string  `json:"finance_code" gorm:"column:finance_code"`
	Total       float64 `json:"total" gorm:"column:total"`
	DeptID      int64   `json:"dept_id" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderAccountSet ERP订单账户保存参数
type OrderAccountSet struct {
	FinanceCode string  `json:"finance_code" jsonschema:"description=财务科目编码"`
	Total       float64 `json:"total" jsonschema:"description=账户金额"`
}
