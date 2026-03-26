package settingServiceImpl

import (
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
