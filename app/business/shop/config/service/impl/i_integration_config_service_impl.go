package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"
	api2 "nova-factory-server/app/utils/observer/integration/api"

	"github.com/gin-gonic/gin"
)

type ShopErpIntegrationConfigServiceImpl struct {
	dao dao.ShopErpIIntegrationConfigDao
}

func NewIntegrationConfigService(dao dao.ShopErpIIntegrationConfigDao) service.ShopErpIIntegrationConfigService {
	return &ShopErpIntegrationConfigServiceImpl{dao: dao}
}

func (i *ShopErpIntegrationConfigServiceImpl) Set(c *gin.Context, req *models.ShopErpIntegrationConfigSet) (*models.ShopErpIntegrationConfig, error) {
	return i.dao.Set(c, req)
}

func (i *ShopErpIntegrationConfigServiceImpl) List(c *gin.Context, req *models.ShopErpIntegrationConfigQuery) (*models.ShopErpIntegrationConfigListData, error) {
	return i.dao.List(c, req)
}

func (i *ShopErpIntegrationConfigServiceImpl) CheckLoginState(c *gin.Context, req *models.ShopErpIntegrationConfigCheckLoginReq) (api2.LoginState, error) {
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

func (i *ShopErpIntegrationConfigServiceImpl) GetEnabled(c *gin.Context) (*models.ShopErpIntegrationConfig, error) {
	cfg, err := i.dao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到启用的集成配置")
	}
	return cfg, nil
}
