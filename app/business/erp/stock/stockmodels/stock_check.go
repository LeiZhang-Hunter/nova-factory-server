package stockmodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// StockCheck ERP 库存盘点单
type StockCheck struct {
	ID         int64      `json:"id,string" gorm:"column:id"`
	No         string     `json:"no" gorm:"column:no"`
	CheckTime  *time.Time `json:"checkTime" gorm:"column:check_time"`
	TotalCount float64    `json:"totalCount" gorm:"column:total_count"`
	TotalPrice float64    `json:"totalPrice" gorm:"column:total_price"`
	Status     int32      `json:"status" gorm:"column:status"`
	Remark     string     `json:"remark" gorm:"column:remark"`
	FileURL    string     `json:"fileUrl" gorm:"column:file_url"`
	DeptID     int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockCheckUpsert ERP 库存盘点单新增修改参数
type StockCheckUpsert struct {
	ID         int64   `json:"id,string"`
	No         string  `json:"no" binding:"required" label:"盘点单号"`
	CheckTime  string  `json:"checkTime"`
	TotalCount float64 `json:"totalCount"`
	TotalPrice float64 `json:"totalPrice"`
	Status     int32   `json:"status"`
	Remark     string  `json:"remark"`
	FileURL    string  `json:"fileUrl"`
}

// StockCheckQuery ERP 库存盘点单查询参数
type StockCheckQuery struct {
	No     string `form:"no" filter:"like,no"`
	Status *int32 `form:"status" filter:"eq,status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

// StockCheckListData ERP 库存盘点单分页数据
type StockCheckListData struct {
	Rows  []*StockCheck `json:"rows"`
	Total int64         `json:"total"`
}
