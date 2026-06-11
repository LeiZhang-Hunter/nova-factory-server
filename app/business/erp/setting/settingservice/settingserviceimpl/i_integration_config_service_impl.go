package settingserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/business/erp/setting/settingservice"
	api2 "nova-factory-server/app/utils/observer/integration/api"

	"github.com/gin-gonic/gin"
)

type IntegrationConfigServiceImpl struct {
	dao settingdao.IIntegrationConfigDao
}

func NewIntegrationConfigService(dao settingdao.IIntegrationConfigDao) settingservice.IIntegrationConfigService {
	return &IntegrationConfigServiceImpl{dao: dao}
}

func (i *IntegrationConfigServiceImpl) Set(c *gin.Context, req *settingmodels.IntegrationConfigSet) (*settingmodels.IntegrationConfig, error) {
	return i.dao.Set(c, req)
}

func (i *IntegrationConfigServiceImpl) List(c *gin.Context, req *settingmodels.IntegrationConfigQuery) (*settingmodels.IntegrationConfigListData, error) {
	return i.dao.List(c, req)
}

func (i *IntegrationConfigServiceImpl) CheckLoginState(c *gin.Context, req *settingmodels.IntegrationConfigCheckLoginReq) (api2.LoginState, error) {
	cfg, err := i.dao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到启用的集成配置")
	}

	service, err := cfg.Service()
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("没有配置集成商")
	}

	return service.CheckLoginState(cfg, req.RedirectURL)
}

func (i *IntegrationConfigServiceImpl) GetEnabled(c *gin.Context) (*settingmodels.IntegrationConfig, error) {
	cfg, err := i.dao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到启用的集成配置")
	}
	return cfg, nil
}
