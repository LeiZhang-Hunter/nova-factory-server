package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopGoodsDao  商品数据访问接口
type IApiShopGoodsDao interface {
	GetByID(c *gin.Context, id int64) (*models.Goods, error)
	GetByGoodsID(c *gin.Context, goodsID string) (*models.Goods, error)
	List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error)
}
