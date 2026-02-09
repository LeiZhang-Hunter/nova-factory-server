package configurationServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/configuration/configurationDao"
	"nova-factory-server/app/business/configuration/configurationModels"
	"nova-factory-server/app/business/configuration/configurationService"
)

type ConfigurationServiceImpl struct {
	dao configurationDao.ConfigurationDao
}

func NewConfigurationServiceImpl(dao configurationDao.ConfigurationDao) configurationService.ConfigurationService {
	return &ConfigurationServiceImpl{
		dao: dao,
	}
}

func (s *ConfigurationServiceImpl) List(c *gin.Context, req *configurationModels.SysConfigurationReq) (*configurationModels.SysConfigurationList, error) {
	return s.dao.List(c, req)
}

func (s *ConfigurationServiceImpl) Set(c *gin.Context, data *configurationModels.SetSysConfiguration) (*configurationModels.SysConfiguration, error) {
	return s.dao.Set(c, data)
}

func (s *ConfigurationServiceImpl) Remove(c *gin.Context, ids []string) error {
	return s.dao.Remove(c, ids)
}
