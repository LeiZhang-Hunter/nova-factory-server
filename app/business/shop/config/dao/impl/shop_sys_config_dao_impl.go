package impl

import (
	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/constant/commonStatus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShopSysConfigDaoImpl 提供商城系统配置表的数据访问能力。
type ShopSysConfigDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopSysConfigDao 创建商城系统配置 DAO。
func NewShopSysConfigDao(ms *gorm.DB) dao.IShopSysConfigDao {
	return &ShopSysConfigDaoImpl{
		db:        ms,
		tableName: "shop_sys_config",
	}
}

// GetByConfigKey 根据配置键名获取配置
func (s *ShopSysConfigDaoImpl) GetByConfigKey(c *gin.Context, configKey string) (*models.ShopSysConfig, error) {
	var result models.ShopSysConfig
	err := s.db.Table(s.tableName).
		Where("config_key = ?", configKey).
		Where("state = ?", commonStatus.NORMAL).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateByConfigKey 根据配置键名更新配置值
func (s *ShopSysConfigDaoImpl) UpdateByConfigKey(c *gin.Context, configKey string, configValue string) error {
	return s.db.Table(s.tableName).
		Where("config_key = ?", configKey).
		Where("state = ?", commonStatus.NORMAL).
		Update("config_value", configValue).Error
}
