package orderDao

import (
	"nova-factory-server/app/business/erp/setting/settingModels"

	"github.com/gin-gonic/gin"
)

type IOrderDao interface {
	GetEnabledGJPCfg(c *gin.Context) (*settingModels.IntegrationConfig, error)
}
