package dao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IApiShopSkuDao 商品 SKU 数据访问接口
type IApiShopSkuDao interface {
	GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error)
	DeductStock(c *gin.Context, id int64, quantity int64) error
}
