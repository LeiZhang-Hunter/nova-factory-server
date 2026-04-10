package shopmodels

import "time"

// Goods 商品信息
type Goods struct {
	ID            uint64    `json:"id" db:"id"`                        // 主键ID
	GoodsID       string    `json:"goodsId" db:"goods_id"`             // 商品业务ID
	GoodsName     string    `json:"goodsName" db:"goods_name"`         // 商品名称
	GoodsCode     string    `json:"goodsCode" db:"goods_code"`         // 商品编码
	OuterID       string    `json:"outerId" db:"outer_id"`             // 外部系统ID
	ImageURL      string    `json:"imageUrl" db:"image_url"`           // 主图地址
	RetailPrice   float64   `json:"retailPrice" db:"retail_price"`     // 零售价
	GalleryImages string    `json:"galleryImages" db:"gallery_images"` // 图集
	VideoURL      string    `json:"videoUrl" db:"video_url"`           // 视频地址
	Description   string    `json:"description" db:"description"`      // 商品描述
	Weight        float64   `json:"weight" db:"weight"`                // 重量
	WeightUnit    string    `json:"weightUnit" db:"weight_unit"`       // 重量单位
	Unit          string    `json:"unit" db:"unit"`                    // 销售单位
	IsOnSale      int32     `json:"isOnSale" db:"is_on_sale"`          // 是否上架
	Quantity      int64     `json:"quantity" db:"quantity"`            // 库存数量
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`         // 创建时间
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`         // 更新时间
}

// GoodsUpsert 商品新增修改参数
type GoodsUpsert struct {
	ID            uint64  `json:"id"`                           // 主键ID
	GoodsID       string  `json:"goodsId" binding:"required"`   // 商品业务ID
	GoodsName     string  `json:"goodsName" binding:"required"` // 商品名称
	GoodsCode     string  `json:"goodsCode"`                    // 商品编码
	OuterID       string  `json:"outerId"`                      // 外部系统ID
	ImageURL      string  `json:"imageUrl"`                     // 主图地址
	RetailPrice   float64 `json:"retailPrice"`                  // 零售价
	GalleryImages string  `json:"galleryImages"`                // 图集
	VideoURL      string  `json:"videoUrl"`                     // 视频地址
	Description   string  `json:"description"`                  // 商品描述
	Weight        float64 `json:"weight"`                       // 重量
	WeightUnit    string  `json:"weightUnit"`                   // 重量单位
	Unit          string  `json:"unit"`                         // 销售单位
	IsOnSale      int32   `json:"isOnSale"`                     // 是否上架
	Quantity      int64   `json:"quantity"`                     // 库存数量
}

// GoodsQuery 商品查询参数
type GoodsQuery struct {
	GoodsName string `form:"goodsName"` // 商品名称
	GoodsCode string `form:"goodsCode"` // 商品编码
	OuterID   string `form:"outerId"`   // 外部系统ID
	IsOnSale  int32  `form:"isOnSale"`  // 是否上架
	Page      int64  `form:"page"`      // 页码
	Size      int64  `form:"size"`      // 每页数量
}

// GoodsListData 商品列表结果
type GoodsListData struct {
	Rows  []*Goods `json:"rows"`  // 数据列表
	Total int64    `json:"total"` // 总数
}

type ExportGoodsList struct {
	Count   int                 `json:"count"`
	Records []ExportGoodsRecord `json:"records"`
}

type ExportGoodsRecord struct {
	ExternalID string         `json:"external_id"`
	Source     string         `json:"source"`
	Entity     string         `json:"entity"`
	Data       map[string]any `json:"data"`
	SyncedAt   time.Time      `json:"synced_at"`
}
