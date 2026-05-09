package service

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IApiShopGoodsService  商品服务接口
type IApiShopGoodsService interface {
	GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error)
	List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error)
	ListRepurchase(c *gin.Context, userID int64, query *models.GoodsQuery) (*models.GoodsListData, error)
}
