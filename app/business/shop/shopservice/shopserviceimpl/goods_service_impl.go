package shopserviceimpl

import (
	"nova-factory-server/app/business/shop/shopdao"
	"nova-factory-server/app/business/shop/shopmodels"
	"nova-factory-server/app/business/shop/shopservice"

	"github.com/gin-gonic/gin"
)

type ShopGoodsServiceImpl struct {
	dao shopdao.IShopGoodsDao
}

func NewShopGoodsService(dao shopdao.IShopGoodsDao) shopservice.IShopGoodsService {
	return &ShopGoodsServiceImpl{dao: dao}
}

func (s *ShopGoodsServiceImpl) Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	return s.dao.Create(c, req)
}

func (s *ShopGoodsServiceImpl) Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	return s.dao.Update(c, req)
}

func (s *ShopGoodsServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopGoodsServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopGoodsServiceImpl) List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error) {
	return s.dao.List(c, req)
}
