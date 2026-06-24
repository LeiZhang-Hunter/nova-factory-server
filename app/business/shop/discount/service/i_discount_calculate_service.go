package discount

import (
	"github.com/gin-gonic/gin"
)

// IDiscountCalculateService 折扣价格计算服务接口
type IDiscountCalculateService interface {
	// CalculateDiscountPrice 计算商品折扣价
	// userID: 用户ID
	// skuID: SKU ID (shop_goods_sku.sku_id)
	// categoryID: 商品分类数据库ID (shop_category.id)
	// originalPrice: 原价
	// 返回: 折扣价, 是否有折扣
	CalculateDiscountPrice(c *gin.Context, userID int64, goodsID int64, skuID int64, categoryID int64, originalPrice float64) (float64, bool)

	// CalculateSkuDiscountPrice 计算SKU的折扣价
	// skus: SKU列表，会直接修改RetailPrice
	// 返回: 修改后的SKU列表, error
	CalculateSkuDiscountPrice(c *gin.Context, userID int64, categoryID int64, skus []*SkuPrice) ([]*SkuPrice, error)

	// BatchCalculateDiscountPrices 批量计算商品折扣价（内部批量查库一次，用于列表场景消除 N+1）
	// 返回: map[goodsID]折扣价，仅包含有折扣的商品
	BatchCalculateDiscountPrices(c *gin.Context, userID int64, goods []*GoodsWithPrice) map[int64]float64
}

// GoodsWithPrice 商品信息（用于批量折扣计算）
type GoodsWithPrice struct {
	GoodsID     int64
	SkuIDs      []int64 // 该商品的所有 SKU ID 列表
	CategoryID  int64
	RetailPrice float64
}

// SkuPrice SKU价格结构（用于折扣计算）
type SkuPrice struct {
	ID          int64   `json:"id,string"`
	SkuID       int64   `json:"skuId"`
	RetailPrice float64 `json:"retailPrice"` // 折扣后的价格
}
