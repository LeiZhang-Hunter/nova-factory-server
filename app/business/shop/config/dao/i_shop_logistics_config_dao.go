package dao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/config/models"
)

// IShopLogisticsConfigDao 物流配置数据访问接口
type IShopLogisticsConfigDao interface {
	Create(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error)
	Update(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.ShopLogisticsConfig, error)
	GetByType(c *gin.Context, typ string) (*models.ShopLogisticsConfig, error)
	List(c *gin.Context, req *models.ShopLogisticsConfigQuery) (*models.ShopLogisticsConfigListData, error)
	// GetEnabled 读取启动配置
	GetEnabled(c *gin.Context) (*models.ShopLogisticsConfig, error)
}
