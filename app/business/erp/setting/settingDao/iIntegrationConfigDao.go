package settingDao

import (
	"nova-factory-server/app/business/erp/setting/settingModels"

	"github.com/gin-gonic/gin"
)

type IIntegrationConfigDao interface {
	Set(c *gin.Context, req *settingModels.IntegrationConfigSet) (*settingModels.IntegrationConfig, error)
	List(c *gin.Context, req *settingModels.IntegrationConfigQuery) (*settingModels.IntegrationConfigListData, error)
	GetEnabledByType(c *gin.Context, tp string) (*settingModels.IntegrationConfig, error)
}
