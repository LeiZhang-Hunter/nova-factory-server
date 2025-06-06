package daemonizeDaoImpl

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentConfigDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIotAgentConfigDaoImpl(db *gorm.DB) daemonizeDao.IotAgentConfigDao {
	return &IotAgentConfigDaoImpl{

		db:        db,
		tableName: "sys_iot_agent_config",
	}
}

// GetByUuid 根据配置uuid获取配置数据
func (i *IotAgentConfigDaoImpl) GetByUuid(ctx context.Context, uuid string) (*daemonizeModels.SysIotAgentConfig, error) {
	var config *daemonizeModels.SysIotAgentConfig
	ret := i.db.Table(i.tableName).Where("uuid = ?", uuid).First(&config)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return config, nil
}

// GetLastedConfig 获取最新的配置
func (i *IotAgentConfigDaoImpl) GetLastedConfig(ctx context.Context) (*daemonizeModels.SysIotAgentConfig, error) {
	var config *daemonizeModels.SysIotAgentConfig
	ret := i.db.Table(i.tableName).First(&config)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return config, nil
}

// GetLastedConfigList 获取最新的配置列表
func (i *IotAgentConfigDaoImpl) GetLastedConfigList(ctx context.Context, count int) ([]*daemonizeModels.SysIotAgentConfig, error) {
	configList := make([]*daemonizeModels.SysIotAgentConfig, 0)
	ret := i.db.Table(i.tableName).Order("create_time desc").Find(&configList)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return configList, nil
}

// GetVersionListByUuidList 根据uuid列表获取version映射关系
func (i *IotAgentConfigDaoImpl) GetVersionListByUuidList(ctx context.Context, uuidList []string) (versionMap map[string]string, err error) {
	versionMap = make(map[string]string)
	if len(uuidList) == 0 {
		return
	}
	configList := make([]*daemonizeModels.SysIotAgentConfig, 0)
	ret := i.db.Table(i.tableName).Where("uuid in (?)", uuidList).
		Select("uuid", "config_version").Find(&configList)
	if ret.Error != nil {
		return nil, ret.Error
	}
	for _, config := range configList {
		versionMap[config.UUID] = config.ConfigVersion
	}
	return
}

// GetByVersion 根据配置版本号获取配置数据
func (i *IotAgentConfigDaoImpl) GetByVersion(ctx context.Context, configVersion string) (config *daemonizeModels.SysIotAgentConfig, err error) {
	ret := i.db.Table(i.tableName).Where("config_version = ?", configVersion).Find(&config)
	if config == nil {
		return nil, errors.New(fmt.Sprintf("agent config version[%v] not exist in db", configVersion))
	}
	if ret.Error != nil {
		zap.L().Error("agent config version[%v] query db error: %v", zap.Stack(configVersion), zap.Error(err))
		return nil, errors.New(fmt.Sprintf("agent config version[%v] not exist in db", configVersion))
	}
	return
}

// Create 保存配置数据
func (i *IotAgentConfigDaoImpl) Create(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (err error) {
	ret := i.db.Table(i.tableName).Create(&daemonizeModels.SysIotAgentConfig{
		UUID:          config.UUID,
		CompanyUUID:   config.CompanyUUID,
		ConfigVersion: config.ConfigVersion,
		Content:       config.Content,
	})
	if ret.Error != nil {
		zap.L().Error("agent config[%+v] update db error: %v", zap.Any("config", config), zap.Error(err))
		return ret.Error
	}
	return
}
