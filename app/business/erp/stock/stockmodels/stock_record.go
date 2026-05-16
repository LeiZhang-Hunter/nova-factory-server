package stockmodels

import (
	"nova-factory-server/app/baize"
)

// StockRecord ERP 产品库存明细
type StockRecord struct {
	ID          int64   `json:"id,string" gorm:"column:id"`
	ProductID   int64   `json:"productId" gorm:"column:product_id"`
	WarehouseID int64   `json:"warehouseId" gorm:"column:warehouse_id"`
	Count       float64 `json:"count" gorm:"column:count"`
	TotalCount  float64 `json:"totalCount" gorm:"column:total_count"`
	BizType     int32   `json:"bizType" gorm:"column:biz_type"`
	BizID       int64   `json:"bizId" gorm:"column:biz_id"`
	BizItemId   int64   `json:"bizItemId" gorm:"column:biz_item_id"`
	BizNo       string  `json:"bizNo" gorm:"column:biz_no"`
	DeptID      int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// StockRecordUpsert ERP 产品库存明细新增修改参数
type StockRecordUpsert struct {
	ID          int64   `json:"id,string"`
	ProductID   int64   `json:"productId"`
	WarehouseID int64   `json:"warehouseId"`
	Count       float64 `json:"count"`
	TotalCount  float64 `json:"totalCount"`
	BizType     int32   `json:"bizType"`
	BizID       int64   `json:"bizId"`
	BizItemId   int64   `json:"bizItemId"`
	BizNo       string  `json:"bizNo"`
}

// StockRecordQuery ERP 产品库存明细查询参数
type StockRecordQuery struct {
	WarehouseID int64  `form:"warehouseId" filter:"eq,warehouse_id"`
	ProductID   int64  `form:"productId" filter:"eq,product_id"`
	BizType     *int32 `form:"bizType" filter:"eq,biz_type"`
	Page        int64  `form:"page"`
	Size        int64  `form:"size"`
}

// StockRecordListData ERP 产品库存明细分页数据
type StockRecordListData struct {
	Rows  []*StockRecord `json:"rows"`
	Total int64          `json:"total"`
}
