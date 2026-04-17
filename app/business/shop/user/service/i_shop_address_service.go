package service

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopAddressService 商城用户地址服务接口
type IShopAddressService interface {
	Set(c *gin.Context, req *models.AddressSetReq) (*models.Address, error)
	GetByID(c *gin.Context, id int64) (*models.Address, error)
	List(c *gin.Context, req *models.AddressQuery) (*models.AddressListData, error)
	Remove(c *gin.Context, ids []int64) error
}
