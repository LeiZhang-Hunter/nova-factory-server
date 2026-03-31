package configurationDao

import (
	"nova-factory-server/app/business/iot/configuration/configurationModels"

	"github.com/gin-gonic/gin"
)

type ConfigurationDao interface {
	List(c *gin.Context, req *configurationModels.SysConfigurationReq) (*configurationModels.SysConfigurationList, error)
	Set(c *gin.Context, data *configurationModels.SetSysConfiguration) (*configurationModels.SysConfiguration, error)
	Remove(c *gin.Context, ids []string) error
}
