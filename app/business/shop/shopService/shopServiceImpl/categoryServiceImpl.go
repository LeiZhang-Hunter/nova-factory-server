package shopServiceImpl

import (
	"nova-factory-server/app/business/shop/shopDao"
	"nova-factory-server/app/business/shop/shopModels"
	"nova-factory-server/app/business/shop/shopService"

	"github.com/gin-gonic/gin"
)

type ShopCategoryServiceImpl struct {
	dao shopDao.IShopCategoryDao
}

func NewShopCategoryService(dao shopDao.IShopCategoryDao) shopService.IShopCategoryService {
	return &ShopCategoryServiceImpl{dao: dao}
}

func (s *ShopCategoryServiceImpl) Create(c *gin.Context, req *shopModels.CategoryUpsert) (*shopModels.Category, error) {
	return s.dao.Create(c, req)
}

func (s *ShopCategoryServiceImpl) Update(c *gin.Context, req *shopModels.CategoryUpsert) (*shopModels.Category, error) {
	return s.dao.Update(c, req)
}

func (s *ShopCategoryServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopCategoryServiceImpl) GetByID(c *gin.Context, id int64) (*shopModels.Category, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopCategoryServiceImpl) List(c *gin.Context, req *shopModels.CategoryQuery) (*shopModels.CategoryListData, error) {
	return s.dao.List(c, req)
}
