package dao

import (
	"nova-factory-server/app/business/shop/config/models"

	"github.com/gin-gonic/gin"
)

// IShopSysConfigDao 商城系统配置数据访问接口
type IShopSysConfigDao interface {
	GetByConfigKeys(c *gin.Context, configKeys []string) ([]models.ShopSysConfig, error)
	BatchUpdate(c *gin.Context, configs []models.KeyValue) error
}
