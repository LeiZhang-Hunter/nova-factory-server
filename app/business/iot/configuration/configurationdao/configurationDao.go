package configurationdao

import (
	"nova-factory-server/app/business/iot/configuration/configurationmodels"

	"github.com/gin-gonic/gin"
)

type ConfigurationDao interface {
	List(c *gin.Context, req *configurationmodels.SysConfigurationReq) (*configurationmodels.SysConfigurationList, error)
	Set(c *gin.Context, data *configurationmodels.SetSysConfiguration) (*configurationmodels.SysConfiguration, error)
	Remove(c *gin.Context, ids []string) error
}
