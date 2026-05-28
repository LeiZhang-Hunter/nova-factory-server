package impl

import (
	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"

	"github.com/gin-gonic/gin"
)

type ShopSysConfigServiceImpl struct {
	dao dao.IShopSysConfigDao
}

func NewShopSysConfigService(dao dao.IShopSysConfigDao) service.IShopSysConfigService {
	return &ShopSysConfigServiceImpl{dao: dao}
}

// GetByConfigKeys 批量查询配置（通用）
func (s *ShopSysConfigServiceImpl) GetByConfigKeys(c *gin.Context, configKeys []string) ([]models.KeyValue, error) {
	rows, err := s.dao.GetByConfigKeys(c, configKeys)
	if err != nil {
		return nil, err
	}
	configs := make([]models.KeyValue, 0, len(configKeys))
	for _, row := range rows {
		configs = append(configs, models.KeyValue{Key: row.ConfigKey, Value: row.ConfigValue})
	}
	return configs, nil
}

// BatchUpdate 全量覆盖写入配置（通用）
func (s *ShopSysConfigServiceImpl) BatchUpdate(c *gin.Context, configs []models.KeyValue) error {
	return s.dao.BatchUpdate(c, configs)
}
