package shopmodels

import "time"

// Category 商品分类
type Category struct {
	ID           uint64    `json:"id" db:"id"`                      // 主键ID
	ParentID     uint64    `json:"parentId" db:"parent_id"`         // 父级分类ID
	AncestorPath string    `json:"ancestorPath" db:"ancestor_path"` // 祖级路径
	Depth        uint32    `json:"depth" db:"depth"`                // 分类层级
	CategoryName string    `json:"categoryName" db:"category_name"` // 分类名称
	CategoryCode string    `json:"categoryCode" db:"category_code"` // 分类编码
	Sort         int32     `json:"sort" db:"sort"`                  // 排序值
	Status       int32     `json:"status" db:"status"`              // 状态
	IsDeleted    int32     `json:"isDeleted" db:"is_deleted"`       // 是否删除
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`       // 创建时间
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`       // 更新时间
}

// CategoryUpsert 商品分类新增修改参数
type CategoryUpsert struct {
	ID           uint64 `json:"id" form:"id"`                                        // 主键ID
	ParentID     uint64 `json:"parentId" form:"parentId"`                            // 父级分类ID
	CategoryName string `json:"categoryName" form:"categoryName" binding:"required"` // 分类名称
	CategoryCode string `json:"categoryCode" form:"categoryCode"`                    // 分类编码
	Sort         int32  `json:"sort" form:"sort"`                                    // 排序值
	Status       int32  `json:"status" form:"status"`                                // 状态
}

// CategoryQuery 商品分类查询参数
type CategoryQuery struct {
	CategoryName string `form:"categoryName"` // 分类名称
	CategoryCode string `form:"categoryCode"` // 分类编码
	Status       int32  `form:"status"`       // 状态
	Page         int64  `form:"page"`         // 页码
	Size         int64  `form:"size"`         // 每页数量
}

// CategoryListData 商品分类列表结果
type CategoryListData struct {
	Rows  []*Category `json:"rows"`  // 数据列表
	Total int64       `json:"total"` // 总数
}
