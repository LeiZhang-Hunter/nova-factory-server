package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IShopSysConfigDao 商城系统配置数据访问接口
type IShopSysConfigDao interface {
	GetByConfigKey(c *gin.Context, configKey string) (*models.ShopSysConfig, error)
	UpdateByConfigKey(c *gin.Context, configKey string, configValue string) error
}
