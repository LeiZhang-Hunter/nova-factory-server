package shopService

import (
	"nova-factory-server/app/business/shop/shopModels"

	"github.com/gin-gonic/gin"
)

type IShopGoodsService interface {
	Create(c *gin.Context, req *shopModels.GoodsUpsert) (*shopModels.Goods, error)
	Update(c *gin.Context, req *shopModels.GoodsUpsert) (*shopModels.Goods, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopModels.Goods, error)
	List(c *gin.Context, req *shopModels.GoodsQuery) (*shopModels.GoodsListData, error)
}
