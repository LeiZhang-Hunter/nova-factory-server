package dao

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopAddressDao 商城用户地址数据访问接口
type IShopAddressDao interface {
	Set(c *gin.Context, req *models.AddressSetReq) (*models.Address, error)
	List(c *gin.Context, req *models.AddressQuery) (*models.AddressListData, error)
	Remove(c *gin.Context, ids []int64) error
}
