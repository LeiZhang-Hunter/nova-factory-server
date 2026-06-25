package dao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IApiShopSkuDao 商品 SKU 数据访问接口
type IApiShopSkuDao interface {
	GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error)
	GetByIDForUpdate(c *gin.Context, id int64) (*shopmodels.GoodsSku, error)
	ListByIDs(c *gin.Context, ids []int64) ([]*shopmodels.GoodsSku, error)
	ListBySkuIDs(c *gin.Context, ids []int64) ([]*shopmodels.GoodsSku, error)
	ListByGoodsIDs(c *gin.Context, goodsIDs []int64) ([]*shopmodels.GoodsSku, error)
	DeductStock(c *gin.Context, id int64, quantity int64) error
	RestoreStock(c *gin.Context, id int64, quantity int64) error
	GetBySkuID(c *gin.Context, skuId int64) (*shopmodels.GoodsSku, error)
	// GetBySkuIDForUpdate 在当前事务中按主键锁定 SKU 行。
	GetBySkuIDForUpdate(c *gin.Context, id int64) (*shopmodels.GoodsSku, error)
	// DeductStockBySkuId 原子扣减 SKU by sku id
	DeductStockBySkuId(c *gin.Context, id int64, quantity int64) error
}
