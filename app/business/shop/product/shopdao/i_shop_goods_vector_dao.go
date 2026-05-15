package shopdao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IShopGoodsVectorDao 商品向量数据访问接口
type IShopGoodsVectorDao interface {
	Upsert(c *gin.Context, goods *shopmodels.Goods, items []*shopmodels.GoodsVectorUpsertItem) (*shopmodels.GoodsVectorResult, error)
	Search(c *gin.Context, req *shopmodels.GoodsVectorSearchReq, vector []float32) (*shopmodels.GoodsVectorSearchData, error)
}
