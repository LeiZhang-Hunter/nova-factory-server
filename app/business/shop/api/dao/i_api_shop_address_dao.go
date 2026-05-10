package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopAddressDao  移动端地址数据访问接口
type IApiShopAddressDao interface {
	Set(c *gin.Context, req *models.AddressSetReq) (*models.ShopUserAddressApp, error)
	GetByID(c *gin.Context, id int64) (*models.ShopUserAddressApp, error)
	List(c *gin.Context, userId int64) (*models.AddressListData, error)
	Remove(c *gin.Context, ids []int64) error
	ClearDefault(c *gin.Context, userId int64, excludeId int64) error
	Default(c *gin.Context, userId int64) (*models.ShopUserAddressApp, error)
}
