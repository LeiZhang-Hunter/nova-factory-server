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

// GetWechatPayConfig 获取微信支付配置
func (s *IApiShopSysConfigDaoImpl) GetWechatPayConfig(c *gin.Context) (*models.ShopSysConfigWechatPayConfigDTO, error) {
	var configKeys []string = []string{
		"wechat_mini_program_app_id",
		"wechat_pay_mch_id",
		"wechat_pay_api_v3_key",
		"wechat_pay_serial_no",
		"wechat_pay_private_key_path",
		"wechat_pay_notify_url",
		"wechat_pay_platform_public_key_id",
		"wechat_pay_platform_public_key_path",
	}

	rows, err := s.GetByConfigKeys(c, configKeys)
	if err != nil {
		return nil, err
	}
	cfgMap := make(map[string]string)
	// 初始化空值
	for _, v := range configKeys {
		cfgMap[v] = ""
	}
	for _, row := range rows {
		cfgMap[row.ConfigKey] = row.ConfigValue
	}
	return &models.ShopSysConfigWechatPayConfigDTO{
		AppId:                 cfgMap["wechat_mini_program_app_id"],
		MchId:                 cfgMap["wechat_pay_mch_id"],
		ApiV3Key:              cfgMap["wechat_pay_api_v3_key"],
		SerialNo:              cfgMap["wechat_pay_serial_no"],
		PrivateKeyPath:        cfgMap["wechat_pay_private_key_path"],
		NotifyUrl:             cfgMap["wechat_pay_notify_url"],
		PlatformPublicKeyId:   cfgMap["wechat_pay_platform_public_key_id"],
		PlatformPublicKeyPath: cfgMap["wechat_pay_platform_public_key_path"],
	}, nil
}

// GetIsAutoRefundEnabled 获取售后订单是否启用自动退款
func (s *IApiShopSysConfigDaoImpl) GetIsAutoRefundEnabled(c *gin.Context) (bool, error) {
	cfg, err := s.GetByConfigKey(c, "shop_auto_refund_unshipped_enabled")
	if err != nil {
		return false, nil // 默认关闭
	}
	return cfg.ConfigValue == "true", nil
}
