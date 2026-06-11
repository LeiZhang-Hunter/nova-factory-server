package models

import "time"

type QQDConfig struct {
	AppKey          string `json:"appKey"`
	AppSecret       string `json:"appSecret"`
	Selfmallaccount string `json:"selfmallaccount"`
	CodeTTL         string `json:"codeTTL"`
	TokenTTL        string `json:"tokenTTL"`
	RefreshTokenTTL string `json:"refreshTokenTTL"`
}

func (c *QQDConfig) ApplyDefaults() {
	if c.CodeTTL == "" {
		c.CodeTTL = "10m"
	}
	if c.TokenTTL == "" {
		c.TokenTTL = "24h"
	}
	if c.RefreshTokenTTL == "" {
		c.RefreshTokenTTL = "720h"
	}
}

type QQDProductTable struct {
	GoodsID        string    `gorm:"column:goods_id"`
	GoodsName      string    `gorm:"column:goods_name"`
	OuterID        string    `gorm:"column:outer_id"`
	ImageURL       string    `gorm:"column:image_url"`
	RetailPrice    float64   `gorm:"column:retail_price"`
	Description    string    `gorm:"column:description"`
	Quantity       int32     `gorm:"column:quantity"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UpdateTime     time.Time `gorm:"column:update_time"`
	ShopCategoryID int64     `gorm:"column:shop_category_id"`
}

type QQDProductSkuTable struct {
	GoodsID     string    `gorm:"column:goods_id"`
	SkuID       string    `gorm:"column:sku_id"`
	SkuName     string    `gorm:"column:sku_name"`
	OuterID     string    `gorm:"column:outer_id"`
	RetailPrice float64   `gorm:"column:retail_price"`
	Quantity    int32     `gorm:"column:quantity"`
	CreateTime  time.Time `gorm:"column:create_time"`
	UpdateTime  time.Time `gorm:"column:update_time"`
}

type QQDSkuStockUpdate struct {
	SkuID    string
	Quantity int64
}
