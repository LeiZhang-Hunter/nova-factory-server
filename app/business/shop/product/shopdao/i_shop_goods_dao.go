package shopdao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IShopGoodsDao interface {
	Transaction(c *gin.Context, fn func(txDao IShopGoodsDao) error) error
	Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)
	BatchCreate(c *gin.Context, reqs []*shopmodels.GoodsUpsert, batchSize int) error
	Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)
	BatchUpdate(c *gin.Context, reqs []*shopmodels.GoodsUpsert, batchSize int) error
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error)
	GetByGoodsID(c *gin.Context, goodsID int64) (*shopmodels.Goods, error)
	ListByGoodsIDs(c *gin.Context, goodsIDs []int64) ([]*shopmodels.Goods, error)
	List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error)
	UpdateStockByGoodsIDWithDB(db *gorm.DB, goodsID int64, quantity int64) error
	UpdateStockByGoodsID(c *gin.Context, goodsID string, quantity int64) error
	UpsertByGoodsID(c *gin.Context, goodsID string, req *shopmodels.GoodsSyncUpsert) error
	UpsertByGoodsIDWithDB(db *gorm.DB, goodsID int64, req *shopmodels.GoodsSyncUpsert) (int64, error)
	LockStockRows(db *gorm.DB, goodsIDs []int64) error
	GetDBInfoByGoodsID(c *gin.Context, goodsID int64) (*shopmodels.Goods, error)
}
