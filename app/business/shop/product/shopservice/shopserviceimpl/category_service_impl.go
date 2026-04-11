package shopserviceimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"

	"github.com/gin-gonic/gin"
)

type ShopCategoryServiceImpl struct {
	dao shopdao.IShopCategoryDao
}

func NewShopCategoryService(dao shopdao.IShopCategoryDao) shopservice.IShopCategoryService {
	return &ShopCategoryServiceImpl{dao: dao}
}

func (s *ShopCategoryServiceImpl) Create(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	return s.dao.Create(c, req)
}

func (s *ShopCategoryServiceImpl) Update(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	return s.dao.Update(c, req)
}

func (s *ShopCategoryServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopCategoryServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Category, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopCategoryServiceImpl) List(c *gin.Context, req *shopmodels.CategoryQuery) (*shopmodels.CategoryListData, error) {
	return s.dao.List(c, req)
}
