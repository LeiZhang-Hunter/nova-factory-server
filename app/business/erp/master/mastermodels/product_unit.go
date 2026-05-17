package mastermodels

import (
	"nova-factory-server/app/baize"
)

// ProductUnit ERP 产品单位
type ProductUnit struct {
	ID     int64  `json:"id,string" gorm:"column:id"`
	Name   string `json:"name" gorm:"column:name"`
	Status int32  `json:"status" gorm:"column:status"`
	DeptID int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// ProductUnitUpsert ERP 产品单位新增修改参数
type ProductUnitUpsert struct {
	ID     int64  `json:"id,string"`
	Name   string `json:"name" binding:"required" label:"单位名字"`
	Status int32  `json:"status"`
}

// ProductUnitQuery ERP 产品单位查询参数
type ProductUnitQuery struct {
	Name   string `form:"name" filter:"like,name"`
	Status *int32 `form:"status" filter:"eq,status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

// ProductUnitListData ERP 产品单位分页数据
type ProductUnitListData struct {
	Rows  []*ProductUnit `json:"rows"`
	Total int64          `json:"total"`
}
