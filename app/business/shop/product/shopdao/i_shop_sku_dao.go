package shopdao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IShopSkuDao interface {
	Transaction(c *gin.Context, fn func(txDao IShopSkuDao) error) error
	Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error)
	BatchCreate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error
	Update(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error)
	BatchUpdate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error)
	ListByIDs(c *gin.Context, ids []int64) ([]*shopmodels.GoodsSku, error)
	GetBySkuID(c *gin.Context, skuID string) (*shopmodels.GoodsSku, error)
	ListByGoodsIDs(c *gin.Context, goodsIDs []int64) ([]*shopmodels.GoodsSku, error)
	ListBySkuIDs(c *gin.Context, skuIDs []int64) ([]*shopmodels.GoodsSku, error)
	List(c *gin.Context, req *shopmodels.GoodsSkuQuery) (*shopmodels.GoodsSkuListData, error)
	UpdateStockBySkuIDWithDB(db *gorm.DB, skuID string, quantity int64) error
	UpdateStockBySkuID(c *gin.Context, skuID string, quantity int64) error
	SumStockByGoodsID(c *gin.Context, goodsID string) (int64, error)
	UpsertBySkuID(c *gin.Context, skuID string, req *shopmodels.GoodsSkuSyncUpsert) error
	UpsertBySkuIDWithDB(db *gorm.DB, skuID string, req *shopmodels.GoodsSkuSyncUpsert) error
	LockStockRows(db *gorm.DB, goodsIDs []int64) error
}
