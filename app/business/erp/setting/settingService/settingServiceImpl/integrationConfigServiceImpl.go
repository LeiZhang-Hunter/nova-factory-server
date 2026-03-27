package settingServiceImpl

import (
	"errors"
	"nova-factory-server/app/business/erp/core/integration"
	"nova-factory-server/app/business/erp/core/integration/api"
	"nova-factory-server/app/business/erp/setting/settingDao"
	"nova-factory-server/app/business/erp/setting/settingModels"
	"nova-factory-server/app/business/erp/setting/settingService"

	"github.com/gin-gonic/gin"
)

type IntegrationConfigServiceImpl struct {
	dao settingDao.IIntegrationConfigDao
}

func NewIntegrationConfigService(dao settingDao.IIntegrationConfigDao) settingService.IIntegrationConfigService {
	return &IntegrationConfigServiceImpl{dao: dao}
}

func (i *IntegrationConfigServiceImpl) Set(c *gin.Context, req *settingModels.IntegrationConfigSet) (*settingModels.IntegrationConfig, error) {
	return i.dao.Set(c, req)
}

func (i *IntegrationConfigServiceImpl) List(c *gin.Context, req *settingModels.IntegrationConfigQuery) (*settingModels.IntegrationConfigListData, error) {
	return i.dao.List(c, req)
}

func (i *IntegrationConfigServiceImpl) CheckLoginState(c *gin.Context, req *settingModels.IntegrationConfigCheckLoginReq) (*api.LoginState, error) {
	cfg, err := i.dao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到启用的集成配置")
	}
	client, err := integration.CreateByType(cfg.Type)
	if err != nil {
		return nil, err
	}
	return client.CheckLoginState(c, cfg, req.CheckURL, req.RedirectURL)
}

func (i *IntegrationConfigServiceImpl) GetEnabled(c *gin.Context) (*settingModels.IntegrationConfig, error) {
	cfg, err := i.dao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到启用的集成配置")
	}
	return cfg, nil
}
