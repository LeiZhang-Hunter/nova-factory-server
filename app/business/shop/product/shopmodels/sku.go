package shopmodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// GoodsSku 商品规格信息
type GoodsSku struct {
	ID                 uint64   `json:"id,string" gorm:"id"`                 // 主键ID
	GoodsDBID          int64    `json:"goodsDBId,string" gorm:"goods_db_id"` // 商品业务ID
	GoodsID            int64    `json:"goodsId,string" gorm:"goods_id"`      // 商品业务ID
	GoodsCode          string   `json:"goodsCode" gorm:"goods_code"`         // 商品编码
	SkuID              int64    `json:"skuId,string" gorm:"sku_id"`          // 规格业务ID
	SkuName            string   `json:"skuName" gorm:"sku_name"`             // 规格名称
	SkuCode            string   `json:"skuCode" gorm:"sku_code"`             // 规格编码
	OuterID            string   `json:"outerId" gorm:"outer_id"`             // 外部系统ID
	Barcode            string   `json:"barcode" gorm:"barcode"`              // 条码
	ImageURL           string   `json:"imageUrl" gorm:"image_url"`           // 主图地址
	RetailPrice        float64  `json:"retailPrice" gorm:"retail_price"`     // 零售价
	GalleryImages      string   `json:"-" gorm:"gallery_images"`             // 图集
	GalleryImagesArray []string `json:"galleryImages" gorm:"-" gorm:"-"`     // 图集
	VideoURL           string   `json:"videoUrl" gorm:"video_url"`           // 视频地址
	Description        string   `json:"description" gorm:"description"`      // 规格描述
	Weight             float64  `json:"weight" gorm:"weight"`                // 重量
	WeightUnit         string   `json:"weightUnit" gorm:"weight_unit"`       // 重量单位
	Unit               string   `json:"unit" gorm:"unit"`                    // 销售单位
	Quantity           int64    `json:"quantity" gorm:"quantity"`            // 库存数量
	DeptID             int64    `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// GoodsSkuUpsert 商品规格新增修改参数
type GoodsSkuUpsert struct {
	ID                 uint64   `json:"id,string"`                           // 主键ID
	GoodsID            int64    `json:"goodsId,string" binding:"required"`   // 商品业务ID
	SkuID              int64    `json:"skuId,string" binding:"required"`     // 规格业务ID
	GoodsDBID          int64    `json:"goodsDBId,string" gorm:"goods_db_id"` // 商品业务ID
	SkuName            string   `json:"skuName"`                             // 规格名称
	SkuCode            string   `json:"skuCode"`                             // 规格编码
	OuterID            string   `json:"outerId"`                             // 外部系统ID
	Barcode            string   `json:"barcode"`                             // 条码
	ImageURL           string   `json:"imageUrl"`                            // 主图地址
	RetailPrice        float64  `json:"retailPrice"`                         // 零售价
	GalleryImagesArray []string `json:"galleryImages" `
	VideoURL           string   `json:"videoUrl"`    // 视频地址
	Description        string   `json:"description"` // 规格描述
	Weight             float64  `json:"weight"`      // 重量
	WeightUnit         string   `json:"weightUnit"`  // 重量单位
	Unit               string   `json:"unit"`        // 销售单位
	Quantity           int64    `json:"quantity"`    // 库存数量
	baize.BaseEntity
}

// GoodsSkuQuery 商品规格查询参数
type GoodsSkuQuery struct {
	GoodsID int64  `form:"goodsId"` // 商品业务ID
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

// GoodsSkuSyncUpsert SKU 同步 upsert 参数，用于 SyncEvent 等场景。
type GoodsSkuSyncUpsert struct {
	GoodsDBId   int64   // 商品数据库id
	GoodsID     int64   // 商品业务ID
	SkuID       string  // 规格业务ID
	SkuName     string  // 规格名称
	SkuCode     string  // 规格编码
	OuterID     string  // 外部系统ID
	Barcode     string  // 条码
	RetailPrice float64 // 零售价
	Weight      float64 // 重量
	WeightUnit  string  // 重量单位
	Quantity    int64   // 库存数量
}

// ToSyncMap 转换为数据库字段映射，用于 GORM Create/Updates，仅包含非空/非零字段
func (r *GoodsSkuSyncUpsert) ToSyncMap(now *time.Time) map[string]any {
	m := map[string]any{"update_time": now}
	if r.GoodsID != 0 {
		m["goods_id"] = r.GoodsID
	}
	if r.SkuID != "" {
		m["sku_id"] = r.SkuID
	}
	if r.SkuName != "" {
		m["sku_name"] = r.SkuName
	}
	if r.SkuCode != "" {
		m["sku_code"] = r.SkuCode
	}
	if r.OuterID != "" {
		m["outer_id"] = r.OuterID
	}
	if r.Barcode != "" {
		m["barcode"] = r.Barcode
	}
	if r.RetailPrice != 0 {
		m["retail_price"] = r.RetailPrice
	}
	if r.Weight != 0 {
		m["weight"] = r.Weight
	}
	if r.WeightUnit != "" {
		m["weight_unit"] = r.WeightUnit
	}
	if r.Quantity != 0 {
		m["quantity"] = r.Quantity
	}
	if r.GoodsDBId != 0 {
		m["goods_db_id"] = r.GoodsDBId
	}
	return m
}
