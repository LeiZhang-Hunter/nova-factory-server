package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IShopAddressService 移动端地址服务接口
type IShopAddressService interface {
	Set(c *gin.Context, req *models.AddressSetReq) (*models.ShopUserAddressApp, error)
	GetByID(c *gin.Context, id int64) (*models.ShopUserAddressApp, error)
	List(c *gin.Context, userId int64) (*models.AddressListData, error)
	Remove(c *gin.Context, ids []int64) error
}
