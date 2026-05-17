package models

// Goods 商品信息
type Goods struct {
	ID               int64    `json:"id,string" gorm:"id"`                    // 主键ID
	GoodsID          string   `json:"goodsId" gorm:"goods_id"`                // 商品业务ID
	ShopCategoryId   int64    `json:"shopCategoryId" gorm:"shop_category_id"` // 商品分类id
	ShopCategoryName string   `json:"shopCategoryName" gorm:"-" gorm:"-"`     // 商品分类名称
	GoodsName        string   `json:"goodsName" gorm:"goods_name"`            // 商品名称
	GoodsCode        string   `json:"goodsCode" gorm:"goods_code"`            // 商品编码
	OuterID          string   `json:"outerId" gorm:"outer_id"`                // 外部系统ID
	ImageURL         string   `json:"imageUrl" gorm:"image_url"`              // 主图地址
	RetailPrice      float64  `json:"retailPrice" gorm:"retail_price"`        // 零售价
	GalleryImages    string   `json:"-" gorm:"gallery_images"`                // 图集
	GalleryImagesArr []string `json:"galleryImages" gorm:"-" gorm:"-"`        // 图集
	VideoURL         string   `json:"videoUrl" gorm:"video_url"`              // 视频地址
	Description      string   `json:"description" gorm:"description"`         // 商品描述
	Weight           float64  `json:"weight" gorm:"weight"`                   // 重量
	WeightUnit       string   `json:"weightUnit" gorm:"weight_unit"`          // 重量单位
	Unit             string   `json:"unit" gorm:"unit"`                       // 销售单位
	IsOnSale         int32    `json:"isOnSale" gorm:"is_on_sale"`             // 是否上架
	Quantity         int64    `json:"quantity" gorm:"quantity"`               // 库存数量
}

// GoodsQuery 商品查询参数
type GoodsQuery struct {
	GoodsName  string `form:"goodsName"`  // 商品名称
	GoodsCode  string `form:"goodsCode"`  // 商品编码
	CategoryId int64  `form:"categoryId"` // 商品分类ID
	SortBy     string `form:"sortBy"`     // 排序字段
	SortOrder  string `form:"sortOrder"`  // 排序方向 asc/desc
	Page       int64  `form:"page"`       // 页码
	Size       int64  `form:"size"`       // 每页数量
}

// GoodsListData 商品列表结果
type GoodsListData struct {
	Rows  []*Goods `json:"rows"`  // 数据列表
	Total int64    `json:"total"` // 总数
}
