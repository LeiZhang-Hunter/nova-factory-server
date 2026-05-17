package mastermodels

import (
	"nova-factory-server/app/baize"
)

// ProductCategory ERP 产品分类
type ProductCategory struct {
	ID         int64  `json:"id,string" gorm:"column:id"`
	ParentID   int64  `json:"parentId,string" gorm:"column:parent_id"`
	ParentName string `json:"parentName" gorm:"-"`
	Name       string `json:"name" gorm:"column:name"`
	Code       string `json:"code" gorm:"column:code"`
	Sort       int32  `json:"sort" gorm:"column:sort"`
	Status     int32  `json:"status" gorm:"column:status"`
	DeptID     int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// ProductCategoryUpsert ERP 产品分类新增修改参数
type ProductCategoryUpsert struct {
	ID       int64  `json:"id,string"`
	ParentID int64  `json:"parentId,string"`
	Name     string `json:"name" binding:"required" label:"分类名称"`
	Code     string `json:"code"`
	Sort     int32  `json:"sort"`
	Status   int32  `json:"status"`
}

// ProductCategoryQuery ERP 产品分类查询参数
type ProductCategoryQuery struct {
	Name     string `form:"name" filter:"like,name"`
	Code     string `form:"code" filter:"like,code"`
	Status   *int32 `form:"status" filter:"eq,status"`
	ParentID int64  `form:"parentId" filter:"eq,parent_id"`
	Page     int64  `form:"page"`
	Size     int64  `form:"size"`
}

// ProductCategoryListData ERP 产品分类分页数据
type ProductCategoryListData struct {
	Rows  []*ProductCategory `json:"rows"`
	Total int64              `json:"total"`
}
