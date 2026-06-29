package service

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/config/models"
)

// IShopLogisticsConfigService 物流配置服务接口
type IShopLogisticsConfigService interface {
	Create(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error)
	Update(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.ShopLogisticsConfig, error)
	GetByType(c *gin.Context, typ string) (*models.ShopLogisticsConfig, error)
	List(c *gin.Context, req *models.ShopLogisticsConfigQuery) (*models.ShopLogisticsConfigListData, error)
}
