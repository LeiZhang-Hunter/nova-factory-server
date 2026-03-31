package shopServiceImpl

import (
	"nova-factory-server/app/business/shop/shopDao"
	"nova-factory-server/app/business/shop/shopModels"
	"nova-factory-server/app/business/shop/shopService"

	"github.com/gin-gonic/gin"
)

type ShopGoodsServiceImpl struct {
	dao shopDao.IShopGoodsDao
}

func NewShopGoodsService(dao shopDao.IShopGoodsDao) shopService.IShopGoodsService {
	return &ShopGoodsServiceImpl{dao: dao}
}

func (s *ShopGoodsServiceImpl) Create(c *gin.Context, req *shopModels.GoodsUpsert) (*shopModels.Goods, error) {
	return s.dao.Create(c, req)
}

func (s *ShopGoodsServiceImpl) Update(c *gin.Context, req *shopModels.GoodsUpsert) (*shopModels.Goods, error) {
	return s.dao.Update(c, req)
}

func (s *ShopGoodsServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopGoodsServiceImpl) GetByID(c *gin.Context, id int64) (*shopModels.Goods, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopGoodsServiceImpl) List(c *gin.Context, req *shopModels.GoodsQuery) (*shopModels.GoodsListData, error) {
	return s.dao.List(c, req)
}
