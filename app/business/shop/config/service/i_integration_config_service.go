package service

import (
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/utils/observer/integration/api"

	"github.com/gin-gonic/gin"
)

type ShopErpIIntegrationConfigService interface {
	Set(c *gin.Context, req *models.ShopErpIntegrationConfigSet) (*models.ShopErpIntegrationConfig, error)
	List(c *gin.Context, req *models.ShopErpIntegrationConfigQuery) (*models.ShopErpIntegrationConfigListData, error)
	CheckLoginState(c *gin.Context, req *models.ShopErpIntegrationConfigCheckLoginReq) (api.LoginState, error)
	GetEnabled(c *gin.Context) (*models.ShopErpIntegrationConfig, error)
}
