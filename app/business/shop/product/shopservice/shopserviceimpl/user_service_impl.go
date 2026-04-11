package shopserviceimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"

	"github.com/gin-gonic/gin"
)

type ShopUserServiceImpl struct {
	dao shopdao.IShopUserDao
}

func NewShopUserService(dao shopdao.IShopUserDao) shopservice.IShopUserService {
	return &ShopUserServiceImpl{dao: dao}
}

func (s *ShopUserServiceImpl) Create(c *gin.Context, req *shopmodels.UserUpsert) (*shopmodels.User, error) {
	return s.dao.Create(c, req)
}

func (s *ShopUserServiceImpl) Update(c *gin.Context, req *shopmodels.UserUpsert) (*shopmodels.User, error) {
	return s.dao.Update(c, req)
}

func (s *ShopUserServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopUserServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.User, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopUserServiceImpl) List(c *gin.Context, req *shopmodels.UserQuery) (*shopmodels.UserListData, error) {
	return s.dao.List(c, req)
}
