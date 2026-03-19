package shopService

import (
	"nova-factory-server/app/business/shop/shopModels"

	"github.com/gin-gonic/gin"
)

type IShopSkuService interface {
	Create(c *gin.Context, req *shopModels.GoodsSkuUpsert) (*shopModels.GoodsSku, error)
	Update(c *gin.Context, req *shopModels.GoodsSkuUpsert) (*shopModels.GoodsSku, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopModels.GoodsSku, error)
	List(c *gin.Context, req *shopModels.GoodsSkuQuery) (*shopModels.GoodsSkuListData, error)
}
