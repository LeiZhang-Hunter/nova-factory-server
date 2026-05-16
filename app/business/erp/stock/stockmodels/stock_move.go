package stockmodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// StockMove ERP 库存调拨单
type StockMove struct {
	ID         int64      `json:"id,string" gorm:"column:id"`
	No         string     `json:"no" gorm:"column:no"`
	MoveTime   *time.Time `json:"moveTime" gorm:"column:move_time"`
	TotalCount float64    `json:"totalCount" gorm:"column:total_count"`
	TotalPrice float64    `json:"totalPrice" gorm:"column:total_price"`
	Status     int32      `json:"status" gorm:"column:status"`
	Remark     string     `json:"remark" gorm:"column:remark"`
	FileURL    string     `json:"fileUrl" gorm:"column:file_url"`
	DeptID     int64      `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockMoveUpsert ERP 库存调拨单新增修改参数
type StockMoveUpsert struct {
	ID         int64   `json:"id,string"`
	No         string  `json:"no" binding:"required" label:"调拨单号"`
	MoveTime   string  `json:"moveTime"`
	TotalCount float64 `json:"totalCount"`
	TotalPrice float64 `json:"totalPrice"`
	Status     int32   `json:"status"`
	Remark     string  `json:"remark"`
	FileURL    string  `json:"fileUrl"`
}

// StockMoveQuery ERP 库存调拨单查询参数
type StockMoveQuery struct {
	No     string `form:"no" filter:"like,no"`
	Status *int32 `form:"status" filter:"eq,status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

// StockMoveListData ERP 库存调拨单分页数据
type StockMoveListData struct {
	Rows  []*StockMove `json:"rows"`
	Total int64        `json:"total"`
}
