// 管家婆全渠道系统配置与数据表结构定义。
// 包含集成配置（QQDConfig）以及商品、SKU 的数据库模型。
package models

import "time"

// QQDConfig 管家婆全渠道集成配置
// 存储在系统集成配置表中，通过 json.Unmarshal 反序列化
type QQDConfig struct {
	AppKey          string `json:"appKey"`
	AppSecret       string `json:"appSecret"`
	Selfmallaccount string `json:"selfmallaccount"`
	CodeTTL         string `json:"codeTTL"`
	TokenTTL        string `json:"tokenTTL"`
	RefreshTokenTTL string `json:"refreshTokenTTL"`
}

// ApplyDefaults 为未配置的 TTL 字段填充默认值
// codeTTL 默认 10 分钟，tokenTTL 默认 24 小时，refreshTokenTTL 默认 720 小时（30天）
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

// QQDProductTable 商品数据表模型，对应数据库中的 goods 表
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

// QQDProductSkuTable 商品 SKU 数据表模型，对应数据库中的 goods_sku 表
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

// QQDSkuStockUpdate SKU 库存更新结构，用于批量更新商品 SKU 的库存数量
type QQDSkuStockUpdate struct {
	SkuID    string
	Quantity int64
}
