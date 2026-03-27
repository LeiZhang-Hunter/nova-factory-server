package settingService

import (
	"nova-factory-server/app/business/erp/core/integration/api"
	"nova-factory-server/app/business/erp/setting/settingModels"

	"github.com/gin-gonic/gin"
)

type IIntegrationConfigService interface {
	Set(c *gin.Context, req *settingModels.IntegrationConfigSet) (*settingModels.IntegrationConfig, error)
	List(c *gin.Context, req *settingModels.IntegrationConfigQuery) (*settingModels.IntegrationConfigListData, error)
	CheckLoginState(c *gin.Context, req *settingModels.IntegrationConfigCheckLoginReq) (*api.LoginState, error)
	GetEnabled(c *gin.Context) (*settingModels.IntegrationConfig, error)
}
