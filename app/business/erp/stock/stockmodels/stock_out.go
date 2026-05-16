package stockmodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// StockOut ERP 其它出库单
type StockOut struct {
	ID         int64      `json:"id,string" gorm:"column:id"`
	No         string     `json:"no" gorm:"column:no"`
	CustomerID int64      `json:"customerId" gorm:"column:customer_id"`
	OutTime    *time.Time `json:"outTime" gorm:"column:out_time"`
	TotalCount float64    `json:"totalCount" gorm:"column:total_count"`
	TotalPrice float64    `json:"totalPrice" gorm:"column:total_price"`
	Status     int32      `json:"status" gorm:"column:status"`
	Remark     string     `json:"remark" gorm:"column:remark"`
	FileURL    string     `json:"fileUrl" gorm:"column:file_url"`
	DeptID     int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockOutUpsert ERP 其它出库单新增修改参数
type StockOutUpsert struct {
	ID         int64   `json:"id,string"`
	No         string  `json:"no" binding:"required" label:"出库单号"`
	CustomerID int64   `json:"customerId"`
	OutTime    string  `json:"outTime"`
	TotalCount float64 `json:"totalCount"`
	TotalPrice float64 `json:"totalPrice"`
	Status     int32   `json:"status"`
	Remark     string  `json:"remark"`
	FileURL    string  `json:"fileUrl"`
}

// StockOutQuery ERP 其它出库单查询参数
type StockOutQuery struct {
	No         string `form:"no" filter:"like,no"`
	Status     *int32 `form:"status" filter:"eq,status"`
	CustomerID int64  `form:"customerId" filter:"eq,customer_id"`
	Page       int64  `form:"page"`
	Size       int64  `form:"size"`
}

// StockOutListData ERP 其它出库单分页数据
type StockOutListData struct {
	Rows  []*StockOut `json:"rows"`
	Total int64       `json:"total"`
}
