package shopmodels

import "time"

type GoodsSku struct {
	ID            uint64    `json:"id" db:"id"`
	GoodsID       string    `json:"goodsId" db:"goods_id"`
	SkuID         string    `json:"skuId" db:"sku_id"`
	SkuName       string    `json:"skuName" db:"sku_name"`
	SkuCode       string    `json:"skuCode" db:"sku_code"`
	OuterID       string    `json:"outerId" db:"outer_id"`
	Barcode       string    `json:"barcode" db:"barcode"`
	ImageURL      string    `json:"imageUrl" db:"image_url"`
	RetailPrice   float64   `json:"retailPrice" db:"retail_price"`
	GalleryImages string    `json:"galleryImages" db:"gallery_images"`
	VideoURL      string    `json:"videoUrl" db:"video_url"`
	Description   string    `json:"description" db:"description"`
	Weight        float64   `json:"weight" db:"weight"`
	WeightUnit    string    `json:"weightUnit" db:"weight_unit"`
	Unit          string    `json:"unit" db:"unit"`
	Quantity      int64     `json:"quantity" db:"quantity"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type GoodsSkuUpsert struct {
	ID            uint64  `json:"id"`
	GoodsID       string  `json:"goodsId" binding:"required"`
	SkuID         string  `json:"skuId" binding:"required"`
	SkuName       string  `json:"skuName"`
	SkuCode       string  `json:"skuCode"`
	OuterID       string  `json:"outerId"`
	Barcode       string  `json:"barcode"`
	ImageURL      string  `json:"imageUrl"`
	RetailPrice   float64 `json:"retailPrice"`
	GalleryImages string  `json:"galleryImages"`
	VideoURL      string  `json:"videoUrl"`
	Description   string  `json:"description"`
	Weight        float64 `json:"weight"`
	WeightUnit    string  `json:"weightUnit"`
	Unit          string  `json:"unit"`
	Quantity      int64   `json:"quantity"`
}

type GoodsSkuQuery struct {
	GoodsID string `form:"goodsId"`
	SkuName string `form:"skuName"`
	SkuCode string `form:"skuCode"`
	OuterID string `form:"outerId"`
	Barcode string `form:"barcode"`
	Page    int64  `form:"page"`
	Size    int64  `form:"size"`
}

type GoodsSkuListData struct {
	Rows  []*GoodsSku `json:"rows"`
	Total int64       `json:"total"`
}
