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

// GetByConfigKeys 批量查询配置
func (s *ShopSysConfigDaoImpl) GetByConfigKeys(c *gin.Context, configKeys []string) ([]models.ShopSysConfig, error) {
	var results []models.ShopSysConfig
	err := s.db.Table(s.tableName).
		Where("config_key IN ?", configKeys).
		Where("state = ?", commonStatus.NORMAL).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// BatchUpdate 全量覆盖写入配置（先删后插）
func (s *ShopSysConfigDaoImpl) BatchUpdate(c *gin.Context, configs []models.KeyValue) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var keys []string
		for _, cfg := range configs {
			keys = append(keys, cfg.Key)
		}
		if err := tx.Table(s.tableName).
			Where("config_key IN ?", keys).
			Delete(&models.ShopSysConfig{}).Error; err != nil {
			return err
		}
		var rows []models.ShopSysConfig
		for _, cfg := range configs {
			rows = append(rows, models.ShopSysConfig{
				ConfigKey:   cfg.Key,
				ConfigValue: cfg.Value,
			})
		}
		if len(rows) > 0 {
			return tx.Table(s.tableName).CreateInBatches(rows, 100).Error
		}
		return nil
	})
}
