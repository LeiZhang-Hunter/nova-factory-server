package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopSeckillService 秒杀服务接口
type IApiShopSeckillService interface {
	// ListConfigs 获取秒杀时间段列表（含当前状态）
	ListConfigs(c *gin.Context) ([]*models.SeckillConfig, error)
	// GetCurrentConfig 获取当前秒杀时段配置
	GetCurrentConfig(c *gin.Context) (*models.SeckillConfig, error)
	// ListGoods 获取秒杀商品列表
	ListGoods(c *gin.Context, query *models.SeckillQuery) (*models.SeckillListData, error)
	// GetGoodsDetail 获取秒杀商品详情
	GetGoodsDetail(c *gin.Context, id int64, timeID int64) (*models.Seckill, error)
}
