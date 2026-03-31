package settingservice

import (
	"nova-factory-server/app/business/erp/core/integration/api"
	"nova-factory-server/app/business/erp/setting/settingmodels"

	"github.com/gin-gonic/gin"
)

type IIntegrationConfigService interface {
	Set(c *gin.Context, req *settingmodels.IntegrationConfigSet) (*settingmodels.IntegrationConfig, error)
	List(c *gin.Context, req *settingmodels.IntegrationConfigQuery) (*settingmodels.IntegrationConfigListData, error)
	CheckLoginState(c *gin.Context, req *settingmodels.IntegrationConfigCheckLoginReq) (*api.LoginState, error)
	GetEnabled(c *gin.Context) (*settingmodels.IntegrationConfig, error)
}
