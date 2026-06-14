package shopmodels

import (
	"nova-factory-server/app/baize"
	store "nova-factory-server/app/utils/store/category"
	goodsstore "nova-factory-server/app/utils/store/goods"
	"sync"
	"time"
)

// Category 商品分类
type Category struct {
	ID           int64  `json:"id,string" gorm:"id"`               // 主键ID
	ParentID     int64  `json:"parentId,string" gorm:"parent_id"`  // 父级分类ID
	AncestorPath string `json:"ancestorPath" gorm:"ancestor_path"` // 祖级路径
	Depth        uint32 `json:"depth" gorm:"depth"`                // 分类层级
	CategoryName string `json:"categoryName" gorm:"category_name"` // 分类名称
	CategoryCode string `json:"categoryCode" gorm:"category_code"` // 分类编码
	ImageURL     string `json:"imageUrl" gorm:"image_url"`         // 分类图片
	Description  string `json:"description" gorm:"description"`    // 分类描述
	Sort         int32  `json:"sort" gorm:"sort"`                  // 排序值
	Status       *bool  `json:"status" gorm:"status"`              // 状态
	DeptID       int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

type CategoryInfo struct {
	ID           int64           `json:"id,string" gorm:"id"`               // 主键ID
	ParentID     int64           `json:"parentId,string" gorm:"parent_id"`  // 父级分类ID
	AncestorPath string          `json:"ancestorPath" gorm:"ancestor_path"` // 祖级路径
	Depth        uint32          `json:"depth" gorm:"depth"`                // 分类层级
	CategoryName string          `json:"categoryName" gorm:"category_name"` // 分类名称
	CategoryCode string          `json:"categoryCode" gorm:"category_code"` // 分类编码
	ImageURL     string          `json:"imageUrl" gorm:"image_url"`         // 分类图片
	Description  string          `json:"description" gorm:"description"`    // 分类描述
	Sort         int32           `json:"sort" gorm:"sort"`                  // 排序值
	Status       *bool           `json:"status" gorm:"status"`              // 状态
	Children     []*CategoryInfo `json:"children" gorm:"-"`
	CreateTime   *time.Time      `json:"createTime" gorm:"create_time"` //创建时间
	UpdateTime   *time.Time      `json:"updateTime" gorm:"update_time"` //修改时间
	mtx          sync.RWMutex    `json:"-"`
}

func (c *CategoryInfo) ToShopCategoryData() store.ShopCategoryData {
	return c
}

func (c *CategoryInfo) Name() string {
	return c.CategoryName
}

func (c *CategoryInfo) CategoryID() int64 {
	return c.ID
}
func (c *CategoryInfo) ChildrenData() []store.ShopCategoryData {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	list := make([]store.ShopCategoryData, 0)
	for _, child := range c.Children {
		list = append(list, child.ToShopCategoryData())
	}
	return list
}
func (c *CategoryInfo) SetChildren(list []store.ShopCategoryData) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	childrens := make([]*CategoryInfo, len(list))
	for k, child := range list {
		childrens[k] = child.(*CategoryInfo)
	}
	c.Children = childrens
	return nil
}

// CategoryUpsert 商品分类新增修改参数
type CategoryUpsert struct {
	ID           int64  `json:"id,string" form:"id"`                                 // 主键ID
	ParentID     int64  `json:"parentId,string" form:"parentId"`                     // 父级分类ID
	CategoryName string `json:"categoryName" form:"categoryName" binding:"required"` // 分类名称
	CategoryCode string `json:"categoryCode" form:"categoryCode"`                    // 分类编码
	ImageURL     string `json:"imageUrl" form:"imageUrl"`                            // 分类图片
	Description  string `json:"description" form:"description"`                      // 分类描述
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

// CategoryData 适配 shopmodels.Category 为 goods.CategoryData 接口，json tag 对齐管家婆 API 标准。
type CategoryData struct {
	Cid             int                       `json:"cid"`
	Name            string                    `json:"name"`
	Parentcid       int                       `json:"parentcid"`
	Sortorder       int                       `json:"sortorder"`
	Isparent        bool                      `json:"isparent"`
	Childcategories []goodsstore.CategoryData `json:"childcategories"`
}

func (c *CategoryData) GetCid() int                                   { return c.Cid }
func (c *CategoryData) GetName() string                               { return c.Name }
func (c *CategoryData) GetParentcid() int                             { return c.Parentcid }
func (c *CategoryData) GetSortorder() int                             { return c.Sortorder }
func (c *CategoryData) GetIsparent() bool                             { return c.Isparent }
func (c *CategoryData) GetChildcategories() []goodsstore.CategoryData { return c.Childcategories }

// CategoryDataResult 实现 goods.DataCategoryResult 接口，json tag 对齐管家婆 API 标准。
type CategoryDataResult struct {
	IsError         bool                      `json:"iserror"`
	ErrorMsg        string                    `json:"errormsg"`
	TotalResults    int                       `json:"-"`
	SellerCategorys []goodsstore.CategoryData `json:"sellercategorys"`
}

func (r *CategoryDataResult) GetIsError() bool                           { return r.IsError }
func (r *CategoryDataResult) GetErrorMsg() string                        { return r.ErrorMsg }
func (r *CategoryDataResult) GetTotalResults() int                       { return r.TotalResults }
func (r *CategoryDataResult) GetCategoryData() []goodsstore.CategoryData { return r.SellerCategorys }

// buildGoodsCategoryTree 递归构建分类树，将 []*shopmodels.Category 转为 goodsstore.CategoryData 接口切片。
func BuildGoodsCategoryTree(childrenMap map[int64][]*Category, parentID int64) []goodsstore.CategoryData {
	children, ok := childrenMap[parentID]
	if !ok || len(children) == 0 {
		return nil
	}
	result := make([]goodsstore.CategoryData, 0, len(children))
	for _, child := range children {
		node := &CategoryData{
			Cid:       int(child.ID),
			Name:      child.CategoryName,
			Parentcid: int(child.ParentID),
			Sortorder: int(child.Sort),
		}
		subChildren := BuildGoodsCategoryTree(childrenMap, child.ID)
		node.Isparent = len(subChildren) > 0
		node.Childcategories = subChildren
		result = append(result, node)
	}
	return result
}
