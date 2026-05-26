package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopSeckillDao 秒杀商品数据访问接口
type IApiShopSeckillDao interface {
	List(c *gin.Context, query *models.SeckillQuery) (*models.SeckillListData, error)
	GetByID(c *gin.Context, id int64) (*models.Seckill, error)
}
