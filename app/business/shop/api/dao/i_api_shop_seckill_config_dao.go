package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopSeckillConfigDao 秒杀配置数据访问接口
type IApiShopSeckillConfigDao interface {
	List(c *gin.Context, query *models.SeckillConfigQuery) (*models.SeckillConfigListData, error)
	GetByID(c *gin.Context, id int64) (*models.SeckillConfig, error)
}
