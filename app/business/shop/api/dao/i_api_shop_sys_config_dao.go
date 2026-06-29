package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopSysConfigDao  商城系统配置数据访问接口
type IApiShopSysConfigDao interface {
	GetByConfigKey(c *gin.Context, configKey string) (*models.ShopSysConfig, error)
	GetByConfigKeys(c *gin.Context, configKeys []string) ([]models.ShopSysConfig, error)
	UpdateByConfigKey(c *gin.Context, configKey string, configValue string) error
	GetWechatPayConfig(c *gin.Context) (*models.ShopSysConfigWechatPayConfigDTO, error)
	GetIsAutoRefundEnabled(c *gin.Context) (bool, error)
}
