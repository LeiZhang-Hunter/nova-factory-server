package shopdao

import (
	"nova-factory-server/app/business/shop/shopmodels"

	"github.com/gin-gonic/gin"
)

type IShopSkuDao interface {
	Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error)
	BatchCreate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error
	Update(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error)
	BatchUpdate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error)
	GetBySkuID(c *gin.Context, skuID string) (*shopmodels.GoodsSku, error)
	ListBySkuIDs(c *gin.Context, skuIDs []string) ([]*shopmodels.GoodsSku, error)
	List(c *gin.Context, req *shopmodels.GoodsSkuQuery) (*shopmodels.GoodsSkuListData, error)
}
