package shopservice

import (
	"nova-factory-server/app/business/shop/shopmodels"

	"github.com/gin-gonic/gin"
)

type IShopGoodsService interface {
	Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)
	Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error)
	List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error)
}
