package shopmodels

import (
	"nova-factory-server/app/baize"
)

// GoodsSku 商品规格信息
type GoodsSku struct {
	ID            uint64  `json:"id" db:"id"`                        // 主键ID
	GoodsID       string  `json:"goodsId" db:"goods_id"`             // 商品业务ID
	SkuID         string  `json:"skuId" db:"sku_id"`                 // 规格业务ID
	SkuName       string  `json:"skuName" db:"sku_name"`             // 规格名称
	SkuCode       string  `json:"skuCode" db:"sku_code"`             // 规格编码
	OuterID       string  `json:"outerId" db:"outer_id"`             // 外部系统ID
	Barcode       string  `json:"barcode" db:"barcode"`              // 条码
	ImageURL      string  `json:"imageUrl" db:"image_url"`           // 主图地址
	RetailPrice   float64 `json:"retailPrice" db:"retail_price"`     // 零售价
	GalleryImages string  `json:"galleryImages" db:"gallery_images"` // 图集
	VideoURL      string  `json:"videoUrl" db:"video_url"`           // 视频地址
	Description   string  `json:"description" db:"description"`      // 规格描述
	Weight        float64 `json:"weight" db:"weight"`                // 重量
	WeightUnit    string  `json:"weightUnit" db:"weight_unit"`       // 重量单位
	Unit          string  `json:"unit" db:"unit"`                    // 销售单位
	Quantity      int64   `json:"quantity" db:"quantity"`            // 库存数量
	DeptID        int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// GoodsSkuUpsert 商品规格新增修改参数
type GoodsSkuUpsert struct {
	ID            uint64  `json:"id"`                         // 主键ID
	GoodsID       string  `json:"goodsId" binding:"required"` // 商品业务ID
	SkuID         string  `json:"skuId" binding:"required"`   // 规格业务ID
	SkuName       string  `json:"skuName"`                    // 规格名称
	SkuCode       string  `json:"skuCode"`                    // 规格编码
	OuterID       string  `json:"outerId"`                    // 外部系统ID
	Barcode       string  `json:"barcode"`                    // 条码
	ImageURL      string  `json:"imageUrl"`                   // 主图地址
	RetailPrice   float64 `json:"retailPrice"`                // 零售价
	GalleryImages string  `json:"galleryImages"`              // 图集
	VideoURL      string  `json:"videoUrl"`                   // 视频地址
	Description   string  `json:"description"`                // 规格描述
	Weight        float64 `json:"weight"`                     // 重量
	WeightUnit    string  `json:"weightUnit"`                 // 重量单位
	Unit          string  `json:"unit"`                       // 销售单位
	Quantity      int64   `json:"quantity"`                   // 库存数量
	baize.BaseEntity
}

// GoodsSkuQuery 商品规格查询参数
type GoodsSkuQuery struct {
	GoodsID string `form:"goodsId"` // 商品业务ID
	SkuName string `form:"skuName"` // 规格名称
	SkuCode string `form:"skuCode"` // 规格编码
	OuterID string `form:"outerId"` // 外部系统ID
	Barcode string `form:"barcode"` // 条码
	Page    int64  `form:"page"`    // 页码
	Size    int64  `form:"size"`    // 每页数量
}

// GoodsSkuListData 商品规格列表结果
type GoodsSkuListData struct {
	Rows  []*GoodsSku `json:"rows"`  // 数据列表
	Total int64       `json:"total"` // 总数
}
