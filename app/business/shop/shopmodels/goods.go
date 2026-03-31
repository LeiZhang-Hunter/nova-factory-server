package shopmodels

import "time"

type Goods struct {
	ID            uint64    `json:"id" db:"id"`
	GoodsID       string    `json:"goodsId" db:"goods_id"`
	GoodsName     string    `json:"goodsName" db:"goods_name"`
	GoodsCode     string    `json:"goodsCode" db:"goods_code"`
	OuterID       string    `json:"outerId" db:"outer_id"`
	ImageURL      string    `json:"imageUrl" db:"image_url"`
	RetailPrice   float64   `json:"retailPrice" db:"retail_price"`
	GalleryImages string    `json:"galleryImages" db:"gallery_images"`
	VideoURL      string    `json:"videoUrl" db:"video_url"`
	Description   string    `json:"description" db:"description"`
	Weight        float64   `json:"weight" db:"weight"`
	WeightUnit    string    `json:"weightUnit" db:"weight_unit"`
	Unit          string    `json:"unit" db:"unit"`
	IsOnSale      int32     `json:"isOnSale" db:"is_on_sale"`
	Quantity      int64     `json:"quantity" db:"quantity"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type GoodsUpsert struct {
	ID            uint64  `json:"id"`
	GoodsID       string  `json:"goodsId" binding:"required"`
	GoodsName     string  `json:"goodsName" binding:"required"`
	GoodsCode     string  `json:"goodsCode"`
	OuterID       string  `json:"outerId"`
	ImageURL      string  `json:"imageUrl"`
	RetailPrice   float64 `json:"retailPrice"`
	GalleryImages string  `json:"galleryImages"`
	VideoURL      string  `json:"videoUrl"`
	Description   string  `json:"description"`
	Weight        float64 `json:"weight"`
	WeightUnit    string  `json:"weightUnit"`
	Unit          string  `json:"unit"`
	IsOnSale      int32   `json:"isOnSale"`
	Quantity      int64   `json:"quantity"`
}

type GoodsQuery struct {
	GoodsName string `form:"goodsName"`
	GoodsCode string `form:"goodsCode"`
	OuterID   string `form:"outerId"`
	IsOnSale  int32  `form:"isOnSale"`
	Page      int64  `form:"page"`
	Size      int64  `form:"size"`
}

type GoodsListData struct {
	Rows  []*Goods `json:"rows"`
	Total int64    `json:"total"`
}
