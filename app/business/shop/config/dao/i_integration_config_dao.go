package dao

import (
	"nova-factory-server/app/business/shop/config/models"

	"github.com/gin-gonic/gin"
)

type ShopErpIIntegrationConfigDao interface {
	Set(c *gin.Context, req *models.ShopErpIntegrationConfigSet) (*models.ShopErpIntegrationConfig, error)
	List(c *gin.Context, req *models.ShopErpIntegrationConfigQuery) (*models.ShopErpIntegrationConfigListData, error)
	GetEnabled(c *gin.Context) (*models.ShopErpIntegrationConfig, error)
}
