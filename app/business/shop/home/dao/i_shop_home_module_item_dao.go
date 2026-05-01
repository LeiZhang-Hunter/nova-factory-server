package dao

import (
	"nova-factory-server/app/business/shop/home/models"

	"github.com/gin-gonic/gin"
)

// IShopHomeModuleItemDao 首页模块明细数据访问接口。
type IShopHomeModuleItemDao interface {
	Set(c *gin.Context, req *models.HomeModuleItemSet) (*models.HomeModuleItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.HomeModuleItem, error)
	List(c *gin.Context, req *models.HomeModuleItemQuery) (*models.HomeModuleItemListData, error)
}
