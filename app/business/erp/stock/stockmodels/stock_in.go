package stockmodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// StockIn ERP 其它入库单
type StockIn struct {
	ID         int64      `json:"id,string" gorm:"column:id"`
	No         string     `json:"no" gorm:"column:no"`
	SupplierID int64      `json:"supplierId" gorm:"column:supplier_id"`
	InTime     *time.Time `json:"inTime" gorm:"column:in_time"`
	TotalCount float64    `json:"totalCount" gorm:"column:total_count"`
	TotalPrice float64    `json:"totalPrice" gorm:"column:total_price"`
	Status     int32      `json:"status" gorm:"column:status"`
	Remark     string     `json:"remark" gorm:"column:remark"`
	FileURL    string     `json:"fileUrl" gorm:"column:file_url"`
	DeptID     int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockInUpsert ERP 其它入库单新增修改参数
type StockInUpsert struct {
	ID         int64   `json:"id,string"`
	No         string  `json:"no" binding:"required" label:"入库单号"`
	SupplierID int64   `json:"supplierId"`
	InTime     string  `json:"inTime"`
	TotalCount float64 `json:"totalCount"`
	TotalPrice float64 `json:"totalPrice"`
	Status     int32   `json:"status"`
	Remark     string  `json:"remark"`
	FileURL    string  `json:"fileUrl"`
}

// StockInQuery ERP 其它入库单查询参数
type StockInQuery struct {
	No         string `form:"no" filter:"like,no"`
	Status     *int32 `form:"status" filter:"eq,status"`
	SupplierID int64  `form:"supplierId" filter:"eq,supplier_id"`
	Page       int64  `form:"page"`
	Size       int64  `form:"size"`
}

// StockInListData ERP 其它入库单分页数据
type StockInListData struct {
	Rows  []*StockIn `json:"rows"`
	Total int64      `json:"total"`
}
