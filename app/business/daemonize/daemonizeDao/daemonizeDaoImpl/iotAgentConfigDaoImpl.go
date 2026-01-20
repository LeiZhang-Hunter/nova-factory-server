package daemonizeDaoImpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
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
func (i *IotAgentConfigDaoImpl) GetLastedConfig(ctx context.Context, agentId uint64) (*daemonizeModels.SysIotAgentConfig, error) {
	var config *daemonizeModels.SysIotAgentConfig
	ret := i.db.Table(i.tableName).Where("agent_object_id = ?", agentId).Order("create_time desc").Where("state = ?", commonStatus.NORMAL).First(&config)
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
func (i *IotAgentConfigDaoImpl) GetVersionListByUuidList(ctx context.Context, uuidList []string) (versionMap map[uint64]string, err error) {
	versionMap = make(map[uint64]string)
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
		versionMap[config.ID] = config.ConfigVersion
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
func (i *IotAgentConfigDaoImpl) Create(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (*daemonizeModels.SysIotAgentConfig, error) {
	if config == nil {
		return nil, errors.New("agent config is nil")
	}
	ret := i.db.Table(i.tableName).Create(config)
	if ret.Error != nil {
		zap.L().Error("agent config[%+v] update db error: %v", zap.Any("config", config), zap.Error(ret.Error))
		return nil, ret.Error
	}
	return config, nil
}

// Update 保存配置数据
func (i *IotAgentConfigDaoImpl) Update(ctx context.Context, config *daemonizeModels.SysIotAgentConfig) (*daemonizeModels.SysIotAgentConfig, error) {
	if config == nil {
		return nil, errors.New("agent config is nil")
	}
	ret := i.db.Table(i.tableName).Where("id = ?", config.ID).Updates(config)
	if ret.Error != nil {
		zap.L().Error("agent config[%+v] update db error: %v", zap.Any("config", config), zap.Error(ret.Error))
		return nil, ret.Error
	}
	return config, nil
}

func (i *IotAgentConfigDaoImpl) List(c *gin.Context, req *daemonizeModels.SysIotAgentConfigListReq) (*daemonizeModels.SysIotAgentConfigListData, error) {
	db := i.db.Table(i.tableName).Table(i.tableName)

	if req != nil && req.AgentObjectID != 0 {
		db = db.Where("agent_object_id = ?", req.AgentObjectID)
	}

	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*daemonizeModels.SysIotAgentConfig

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &daemonizeModels.SysIotAgentConfigListData{
			Rows:  make([]*daemonizeModels.SysIotAgentConfig, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &daemonizeModels.SysIotAgentConfigListData{
			Rows:  make([]*daemonizeModels.SysIotAgentConfig, 0),
			Total: 0,
		}, ret.Error
	}
	return &daemonizeModels.SysIotAgentConfigListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IotAgentConfigDaoImpl) Remove(ctx context.Context, ids []string) error {
	ret := i.db.Table(i.tableName).Where("id in (?)", ids).Delete(&daemonizeModels.SysIotAgentConfig{})
	return ret.Error
}
