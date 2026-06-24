package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopGoodsDao  商品数据访问接口
type IApiShopGoodsDao interface {
	GetByID(c *gin.Context, id int64) (*models.Goods, error)
	ListByIDs(c *gin.Context, ids []int64) ([]*models.Goods, error)
	GetByGoodsID(c *gin.Context, goodsID int64) (*models.Goods, error)
	List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error)
	RandomSale(c *gin.Context, limit int64) (*models.GoodsListData, error)
	ListByUserPurchased(c *gin.Context, userID int64, query *models.GoodsQuery) (*models.GoodsListData, error)
}
