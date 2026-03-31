package shopServiceImpl

import (
	"nova-factory-server/app/business/shop/shopDao"
	"nova-factory-server/app/business/shop/shopModels"
	"nova-factory-server/app/business/shop/shopService"

	"github.com/gin-gonic/gin"
)

type ShopUserServiceImpl struct {
	dao shopDao.IShopUserDao
}

func NewShopUserService(dao shopDao.IShopUserDao) shopService.IShopUserService {
	return &ShopUserServiceImpl{dao: dao}
}

func (s *ShopUserServiceImpl) Create(c *gin.Context, req *shopModels.UserUpsert) (*shopModels.User, error) {
	return s.dao.Create(c, req)
}

func (s *ShopUserServiceImpl) Update(c *gin.Context, req *shopModels.UserUpsert) (*shopModels.User, error) {
	return s.dao.Update(c, req)
}

func (s *ShopUserServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopUserServiceImpl) GetByID(c *gin.Context, id int64) (*shopModels.User, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopUserServiceImpl) List(c *gin.Context, req *shopModels.UserQuery) (*shopModels.UserListData, error) {
	return s.dao.List(c, req)
}
