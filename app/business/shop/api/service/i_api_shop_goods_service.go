package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopGoodsService  商品服务接口
type IApiShopGoodsService interface {
	GetByID(c *gin.Context, id int64) (*models.Goods, error)
	List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error)
}
