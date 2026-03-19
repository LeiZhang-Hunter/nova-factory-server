package shopModels

import "time"

type Category struct {
	ID           uint64    `json:"id" db:"id"`
	ParentID     uint64    `json:"parentId" db:"parent_id"`
	AncestorPath string    `json:"ancestorPath" db:"ancestor_path"`
	Depth        uint32    `json:"depth" db:"depth"`
	CategoryName string    `json:"categoryName" db:"category_name"`
	CategoryCode string    `json:"categoryCode" db:"category_code"`
	Sort         int32     `json:"sort" db:"sort"`
	Status       int32     `json:"status" db:"status"`
	IsDeleted    int32     `json:"isDeleted" db:"is_deleted"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type CategoryUpsert struct {
	ID           uint64 `json:"id" form:"id"`
	ParentID     uint64 `json:"parentId" form:"parentId"`
	CategoryName string `json:"categoryName" form:"categoryName" binding:"required"`
	CategoryCode string `json:"categoryCode" form:"categoryCode"`
	Sort         int32  `json:"sort" form:"sort"`
	Status       int32  `json:"status" form:"status"`
}

type CategoryQuery struct {
	CategoryName string `form:"categoryName"`
	CategoryCode string `form:"categoryCode"`
	Status       int32  `form:"status"`
	Page         int64  `form:"page"`
	Size         int64  `form:"size"`
}

type CategoryListData struct {
	Rows  []*Category `json:"rows"`
	Total int64       `json:"total"`
}
