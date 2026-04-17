package impl

import (
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"

	"github.com/gin-gonic/gin"
)

// ShopAddressServiceImpl 提供商城用户地址相关业务能力。
type ShopAddressServiceImpl struct {
	dao dao.IShopAddressDao
}

// NewShopAddressService 创建商城用户地址服务。
func NewShopAddressService(dao dao.IShopAddressDao) service.IShopAddressService {
	return &ShopAddressServiceImpl{dao: dao}
}

// Set 新增或修改商城用户地址。
func (s *ShopAddressServiceImpl) Set(c *gin.Context, req *models.AddressSetReq) (*models.Address, error) {
	return s.dao.Set(c, req)
}

// List 查询商城用户地址列表。
func (s *ShopAddressServiceImpl) List(c *gin.Context, req *models.AddressQuery) (*models.AddressListData, error) {
	return s.dao.List(c, req)
}

// Remove 删除商城用户地址。
func (s *ShopAddressServiceImpl) Remove(c *gin.Context, ids []int64) error {
	return s.dao.Remove(c, ids)
}
