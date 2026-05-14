package shopdao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IShopGoodsVectorDao 商品向量数据访问接口
type IShopGoodsVectorDao interface {
	Upsert(c *gin.Context, goods *shopmodels.Goods, content string, vector []float32) (*shopmodels.GoodsVectorResult, error)
}
