package shopServiceImpl

import (
	"nova-factory-server/app/business/shop/shopDao"
	"nova-factory-server/app/business/shop/shopModels"
	"nova-factory-server/app/business/shop/shopService"

	"github.com/gin-gonic/gin"
)

type ShopSkuServiceImpl struct {
	dao shopDao.IShopSkuDao
}

func NewShopSkuService(dao shopDao.IShopSkuDao) shopService.IShopSkuService {
	return &ShopSkuServiceImpl{dao: dao}
}

func (s *ShopSkuServiceImpl) Create(c *gin.Context, req *shopModels.GoodsSkuUpsert) (*shopModels.GoodsSku, error) {
	return s.dao.Create(c, req)
}

func (s *ShopSkuServiceImpl) Update(c *gin.Context, req *shopModels.GoodsSkuUpsert) (*shopModels.GoodsSku, error) {
	return s.dao.Update(c, req)
}

func (s *ShopSkuServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopSkuServiceImpl) GetByID(c *gin.Context, id int64) (*shopModels.GoodsSku, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopSkuServiceImpl) List(c *gin.Context, req *shopModels.GoodsSkuQuery) (*shopModels.GoodsSkuListData, error) {
	return s.dao.List(c, req)
}
