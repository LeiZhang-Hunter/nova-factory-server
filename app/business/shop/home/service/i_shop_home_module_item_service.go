package service

import (
	"nova-factory-server/app/business/shop/home/models"

	"github.com/gin-gonic/gin"
)

// IShopHomeModuleItemService 首页模块明细服务接口。
type IShopHomeModuleItemService interface {
	Set(c *gin.Context, req *models.HomeModuleItemSet) (*models.HomeModuleItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.HomeModuleItem, error)
	ListByModuleIDs(c *gin.Context, moduleIDs []int64) ([]*models.HomeModuleItem, error)
	List(c *gin.Context, req *models.HomeModuleItemQuery) (*models.HomeModuleItemListData, error)
}
