package configurationServiceImpl

import (
	"nova-factory-server/app/business/iot/configuration/configurationdao"
	"nova-factory-server/app/business/iot/configuration/configurationmodels"
	"nova-factory-server/app/business/iot/configuration/configurationservice"

	"github.com/gin-gonic/gin"
)

type ConfigurationServiceImpl struct {
	dao configurationdao.ConfigurationDao
}

func NewConfigurationServiceImpl(dao configurationdao.ConfigurationDao) configurationservice.ConfigurationService {
	return &ConfigurationServiceImpl{
		dao: dao,
	}
}

func (s *ConfigurationServiceImpl) List(c *gin.Context, req *configurationmodels.SysConfigurationReq) (*configurationmodels.SysConfigurationList, error) {
	return s.dao.List(c, req)
}

func (s *ConfigurationServiceImpl) Set(c *gin.Context, data *configurationmodels.SetSysConfiguration) (*configurationmodels.SysConfiguration, error) {
	return s.dao.Set(c, data)
}

func (s *ConfigurationServiceImpl) Remove(c *gin.Context, ids []string) error {
	return s.dao.Remove(c, ids)
}
