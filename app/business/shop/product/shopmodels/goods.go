package shopmodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// Goods 商品信息
type Goods struct {
	ID            int64       `json:"id,string" db:"id"`                 // 主键ID
	GoodsID       string      `json:"goodsId" db:"goods_id"`             // 商品业务ID
	GoodsName     string      `json:"goodsName" db:"goods_name"`         // 商品名称
	GoodsCode     string      `json:"goodsCode" db:"goods_code"`         // 商品编码
	OuterID       string      `json:"outerId" db:"outer_id"`             // 外部系统ID
	ImageURL      string      `json:"imageUrl" db:"image_url"`           // 主图地址
	RetailPrice   float64     `json:"retailPrice" db:"retail_price"`     // 零售价
	GalleryImages string      `json:"galleryImages" db:"gallery_images"` // 图集
	VideoURL      string      `json:"videoUrl" db:"video_url"`           // 视频地址
	Description   string      `json:"description" db:"description"`      // 商品描述
	Weight        float64     `json:"weight" db:"weight"`                // 重量
	WeightUnit    string      `json:"weightUnit" db:"weight_unit"`       // 重量单位
	Unit          string      `json:"unit" db:"unit"`                    // 销售单位
	IsOnSale      int32       `json:"isOnSale" db:"is_on_sale"`          // 是否上架
	Quantity      int64       `json:"quantity" db:"quantity"`            // 库存数量
	Skus          []*GoodsSku `json:"skus" gorm:"-" db:"-"`              // 商品规格列表
	DeptID        int64       `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// GoodsUpsert 商品新增修改参数
type GoodsUpsert struct {
	ID            int64   `json:"id,string"`                    // 主键ID
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
	baize.BaseEntity
}

// GoodsQuery 商品查询参数
type GoodsQuery struct {
	GoodsName string `form:"goodsName"` // 商品名称
	GoodsCode string `form:"goodsCode"` // 商品编码
	OuterID   string `form:"outerId"`   // 外部系统ID
	IsOnSale  *bool  `form:"isOnSale"`  // 是否上架
	Page      int64  `form:"page"`      // 页码
	Size      int64  `form:"size"`      // 每页数量
}

// GoodsListData 商品列表结果
type GoodsListData struct {
	Rows  []*Goods `json:"rows"`  // 数据列表
	Total int64    `json:"total"` // 总数
}

// ImportGoodsList 商品导入列表
type ImportGoodsList struct {
	Count   int                 `json:"count"`
	Records []ImportGoodsRecord `json:"records"`
}

// ImportGoodsRecord 导入商品结果
type ImportGoodsRecord struct {
	ExternalID string             `json:"external_id"`
	Source     string             `json:"source"`
	Entity     string             `json:"entity"`
	Data       ImportGoodsRawData `json:"data"`
	SyncedAt   time.Time          `json:"synced_at"`
}

// ImportGoodsSkuRawData 导入的sku原始数据
type ImportGoodsSkuRawData struct {
	Barcode  string  `json:"barcode"`
	Lcmccode string  `json:"lcmccode"`
	Price    float64 `json:"price"`
	Price2   float64 `json:"price2"`
	Price3   float64 `json:"price3"`
	Price4   float64 `json:"price4"`
	Price5   float64 `json:"price5"`
	Size     float64 `json:"size"`
	Skucode  string  `json:"skucode"`
	Skuid    string  `json:"skuid"`
	Skuname  string  `json:"skuname"`
	Weight   float64 `json:"weight"`
}

// ImportGoodsRawData 导入的商品原始数据
type ImportGoodsRawData struct {
	ProductCode string                  `json:"product_code"`
	ProductName string                  `json:"product_name"`
	Remark      string                  `json:"remark"`
	Skus        []ImportGoodsSkuRawData `json:"skus"`
	UnitName    string                  `json:"unit_name"`
	Units       []Unit                  `json:"units"`
}

type Unit struct {
	Unitname string  `json:"unitname"`
	Barcode  string  `json:"barcode"`
	Rate     float64 `json:"rate"`
	Price    float64 `json:"price"`
	Price2   float64 `json:"price2"`
	Price3   float64 `json:"price3"`
	Price4   float64 `json:"price4"`
	Price5   float64 `json:"price5"`
}
