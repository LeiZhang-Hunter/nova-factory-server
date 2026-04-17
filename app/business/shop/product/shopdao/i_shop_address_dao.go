package shopdao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IShopAddressDao 商城地址数据访问接口
type IShopAddressDao interface {
	Set(c *gin.Context, req *shopmodels.AddressSetReq) (*shopmodels.Address, error)
	List(c *gin.Context, req *shopmodels.AddressQuery) (*shopmodels.AddressListData, error)
	Remove(c *gin.Context, ids []int64) error
}
