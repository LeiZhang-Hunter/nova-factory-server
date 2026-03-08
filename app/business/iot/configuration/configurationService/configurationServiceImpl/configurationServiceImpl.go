package configurationServiceImpl

import (
	"nova-factory-server/app/business/iot/configuration/configurationDao"
	"nova-factory-server/app/business/iot/configuration/configurationModels"
	"nova-factory-server/app/business/iot/configuration/configurationService"

	"github.com/gin-gonic/gin"
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
