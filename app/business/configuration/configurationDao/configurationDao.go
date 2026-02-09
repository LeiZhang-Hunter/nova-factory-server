package configurationDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/configuration/configurationModels"
)

type ConfigurationDao interface {
	List(c *gin.Context, req *configurationModels.SysConfigurationReq) (*configurationModels.SysConfigurationList, error)
	Set(c *gin.Context, data *configurationModels.SetSysConfiguration) (*configurationModels.SysConfiguration, error)
	Remove(c *gin.Context, ids []string) error
}
