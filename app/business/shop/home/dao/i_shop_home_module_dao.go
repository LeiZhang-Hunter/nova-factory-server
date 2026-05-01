package dao

import (
	"nova-factory-server/app/business/shop/home/models"

	"github.com/gin-gonic/gin"
)

// IShopHomeModuleDao 首页模块数据访问接口。
type IShopHomeModuleDao interface {
	Set(c *gin.Context, req *models.HomeModuleSet) (*models.HomeModule, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.HomeModule, error)
	List(c *gin.Context, req *models.HomeModuleQuery) (*models.HomeModuleListData, error)
}
