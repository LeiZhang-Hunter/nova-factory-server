package shopmodels

import (
	"nova-factory-server/app/baize"
)

// Category 商品分类
type Category struct {
	ID           int64  `json:"id,string" db:"id"`               // 主键ID
	ParentID     int64  `json:"parentId,string" db:"parent_id"`  // 父级分类ID
	AncestorPath string `json:"ancestorPath" db:"ancestor_path"` // 祖级路径
	Depth        uint32 `json:"depth" db:"depth"`                // 分类层级
	CategoryName string `json:"categoryName" db:"category_name"` // 分类名称
	CategoryCode string `json:"categoryCode" db:"category_code"` // 分类编码
	Sort         int32  `json:"sort" db:"sort"`                  // 排序值
	Status       *bool  `json:"status" db:"status"`              // 状态
	DeptID       int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// CategoryUpsert 商品分类新增修改参数
type CategoryUpsert struct {
	ID           int64  `json:"id,string" form:"id"`                                 // 主键ID
	ParentID     int64  `json:"parentId,string" form:"parentId"`                     // 父级分类ID
	CategoryName string `json:"categoryName" form:"categoryName" binding:"required"` // 分类名称
	CategoryCode string `json:"categoryCode" form:"categoryCode"`                    // 分类编码
	Sort         int32  `json:"sort" form:"sort"`                                    // 排序值
	Status       *bool  `json:"status" form:"status"`                                // 状态
}

// CategoryQuery 商品分类查询参数
type CategoryQuery struct {
	CategoryName string `form:"categoryName"` // 分类名称
	CategoryCode string `form:"categoryCode"` // 分类编码
	Status       *bool  `form:"status"`       // 状态
	Page         int64  `form:"page"`         // 页码
	Size         int64  `form:"size"`         // 每页数量
}

// CategoryListData 商品分类列表结果
type CategoryListData struct {
	Rows  []*Category `json:"rows"`  // 数据列表
	Total int64       `json:"total"` // 总数
}
