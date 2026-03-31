package settingdao

import (
	"nova-factory-server/app/business/erp/setting/settingmodels"

	"github.com/gin-gonic/gin"
)

type IIntegrationConfigDao interface {
	Set(c *gin.Context, req *settingmodels.IntegrationConfigSet) (*settingmodels.IntegrationConfig, error)
	List(c *gin.Context, req *settingmodels.IntegrationConfigQuery) (*settingmodels.IntegrationConfigListData, error)
	GetEnabled(c *gin.Context) (*settingmodels.IntegrationConfig, error)
}
