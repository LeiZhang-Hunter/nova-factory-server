package service

import (
	"nova-factory-server/app/business/shop/config/models"

	"github.com/gin-gonic/gin"
)

// IShopSysConfigService 商城系统配置服务接口（通用）
type IShopSysConfigService interface {
	// GetByConfigKeys 根据 configKey 列表批量获取配置
	GetByConfigKeys(c *gin.Context, configKeys []string) ([]models.KeyValue, error)
	// BatchUpdate 全量覆盖写入配置列表
	BatchUpdate(c *gin.Context, configs []models.KeyValue) error
}
