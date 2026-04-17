package shopserviceimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"

	"github.com/gin-gonic/gin"
)

// ShopAddressServiceImpl 提供商城用户地址相关业务能力。
type ShopAddressServiceImpl struct {
	dao shopdao.IShopAddressDao
}

// NewShopAddressService 创建商城地址服务。
func NewShopAddressService(dao shopdao.IShopAddressDao) shopservice.IShopAddressService {
	return &ShopAddressServiceImpl{dao: dao}
}

// Set 新增或修改商城用户地址。
func (s *ShopAddressServiceImpl) Set(c *gin.Context, req *shopmodels.AddressSetReq) (*shopmodels.Address, error) {
	return s.dao.Set(c, req)
}

// List 查询商城用户地址列表。
func (s *ShopAddressServiceImpl) List(c *gin.Context, req *shopmodels.AddressQuery) (*shopmodels.AddressListData, error) {
	return s.dao.List(c, req)
}

// Remove 删除商城用户地址。
func (s *ShopAddressServiceImpl) Remove(c *gin.Context, ids []int64) error {
	return s.dao.Remove(c, ids)
}
