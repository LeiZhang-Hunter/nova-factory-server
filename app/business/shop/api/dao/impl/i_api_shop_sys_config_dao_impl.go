package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/constant/commonStatus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopSysConfigDaoImpl 提供商城系统配置表的数据访问能力。
type IApiShopSysConfigDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopSysConfigDaoImpl   创建商城系统配置 DAO。
func NewIApiShopSysConfigDaoImpl(ms *gorm.DB) dao.IApiShopSysConfigDao {
	return &IApiShopSysConfigDaoImpl{
		db:        ms,
		tableName: "shop_sys_config",
	}
}

// GetByConfigKeys 批量查询配置
func (s *IApiShopSysConfigDaoImpl) GetByConfigKeys(c *gin.Context, configKeys []string) ([]models.ShopSysConfig, error) {
	var results []models.ShopSysConfig
	err := s.db.WithContext(c).Table(s.tableName).
		Where("config_key IN ?", configKeys).
		Where("state = ?", commonStatus.NORMAL).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// GetByConfigKey 根据配置键名获取配置
func (s *IApiShopSysConfigDaoImpl) GetByConfigKey(c *gin.Context, configKey string) (*models.ShopSysConfig, error) {
	var result models.ShopSysConfig
	err := s.db.WithContext(c).Table(s.tableName).
		Where("config_key = ?", configKey).
		Where("state = ?", commonStatus.NORMAL).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateByConfigKey 根据配置键名更新配置值
func (s *IApiShopSysConfigDaoImpl) UpdateByConfigKey(c *gin.Context, configKey string, configValue string) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("config_key = ?", configKey).
		Where("state = ?", commonStatus.NORMAL).
		Update("config_value", configValue).Error
}
