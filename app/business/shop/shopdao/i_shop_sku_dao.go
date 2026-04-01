package shopdao

import (
	"nova-factory-server/app/business/shop/shopmodels"

	"github.com/gin-gonic/gin"
)

type IShopSkuDao interface {
	Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error)
	Update(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error)
	List(c *gin.Context, req *shopmodels.GoodsSkuQuery) (*shopmodels.GoodsSkuListData, error)
}
