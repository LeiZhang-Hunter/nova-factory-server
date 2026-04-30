package service

import (
	"nova-factory-server/app/business/shop/activity/models"

	"github.com/gin-gonic/gin"
)

// IShopSeckillConfigService 秒杀配置服务接口。
type IShopSeckillConfigService interface {
	Set(c *gin.Context, req *models.SeckillConfigSet) (*models.SeckillConfig, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.SeckillConfig, error)
	GetByIDs(c *gin.Context, ids []int64) ([]*models.SeckillConfig, error)
	List(c *gin.Context, req *models.SeckillConfigQuery) (*models.SeckillConfigListData, error)
}
