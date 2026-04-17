package shopservice

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IShopAddressService 商城地址服务接口
type IShopAddressService interface {
	Set(c *gin.Context, req *shopmodels.AddressSetReq) (*shopmodels.Address, error)
	List(c *gin.Context, req *shopmodels.AddressQuery) (*shopmodels.AddressListData, error)
	Remove(c *gin.Context, ids []int64) error
}
